package openshift

import (
	"time"

	"github.com/openshift/origin/pkg/cmd/cli/cmd"
	"github.com/openshift/origin/pkg/controller"

	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/cache"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/runtime"
	kutil "k8s.io/kubernetes/pkg/util"
	"k8s.io/kubernetes/pkg/watch"

	"github.com/golang/glog"
)

// ExportController represents a controller that will react to changes
// for the given Kind, and push the "exported" resources to a channel.
type ExportController struct {
	// Kind is an instance of the object that should be watched for
	Kind runtime.Object

	// ResourcesChan is the channel over which resources will be send
	ResourcesChan chan<- Resource

	// ResyncPeriod is the interval of time at which the controller
	// will perform of full resync (list) for the provided kind
	ResyncPeriod time.Duration

	// ListFunc is the function used to list the resources for the provided kind
	ListFunc func(options kapi.ListOptions) (runtime.Object, error)

	// WatchFunc is the function used to watch the resources for the provided kind
	WatchFunc func(options kapi.ListOptions) (watch.Interface, error)

	// KeyListFunc is a function that returns the list of keys ("namespace/name" format)
	// that we "know about" (to get a 2-way sync)
	KeyListFunc func() []string

	// KeyGetFunc is a function that returns the object that we "know about"
	// for the given key ("namespace/name" format) - and a boolean if it exists
	KeyGetFunc func(key string) (interface{}, bool, error)

	// Requirements is a slice of functions that can create requirements
	// for customizing the labelSelector for the provided kind
	Requirements []func() (*labels.Requirement, error)

	// LabelSelector is a user-provided labelSelector as a string
	// used to restrict the resources for the provided kind
	LabelSelector string
}

// RunUntil runs the controller in a goroutine
// until stopChan is closed
func (c *ExportController) RunUntil(stopChan <-chan struct{}) {
	queue := cache.NewDeltaFIFO(cache.MetaNamespaceKeyFunc, nil, c)
	cache.NewReflector(c, c.Kind, queue, c.ResyncPeriod).RunUntil(stopChan)

	retryController := &controller.RetryController{
		Handle: c.handle,
		Queue:  queue,
		RetryManager: controller.NewQueueRetryManager(
			queue,
			cache.MetaNamespaceKeyFunc,
			c.retry,
			kutil.NewTokenBucketRateLimiter(1, 10)),
	}

	retryController.RunUntil(stopChan)
}

// handle handles an event change
// by converting it to a known resource
// and pushing it to the output resources channel
func (c *ExportController) handle(obj interface{}) error {
	exporter := cmd.NewExporter()
	deltas := obj.(cache.Deltas)
	for _, delta := range deltas {

		if object, ok := delta.Object.(runtime.Object); ok {
			glog.V(5).Infof("Handling %v for %T", delta.Type, delta.Object)

			// get the reference before exporting
			// (after it will be too late to get a reference)
			ref, err := kapi.GetReference(object)
			if err != nil {
				return err
			}

			if err := exporter.Export(object, false); err != nil {
				if err == cmd.ErrExportOmit {
					// let's just ignore this object that can't be exported
					glog.V(6).Infof("Ignoring %s %s/%s: %v", ref.Kind, ref.Namespace, ref.Name, err)
					continue
				}
				return err
			}

			var exists bool
			switch delta.Type {
			case cache.Added, cache.Updated, cache.Sync:
				exists = true
			case cache.Deleted:
				exists = false
			}

			r := Resource{
				ObjectReference: ref,
				Object:          object,
				Exists:          exists,
				Status:          string(delta.Type),
			}

			glog.V(4).Infof("Processing %s", r)
			c.ResourcesChan <- r

			continue
		}

		if deletedObject, ok := delta.Object.(cache.DeletedFinalStateUnknown); ok {
			glog.V(5).Infof("Handling %v DeletedFinalStateUnknown for %s: %+v", delta.Type, deletedObject.Key, deletedObject.Obj)

			if resource, ok := deletedObject.Obj.(Resource); ok {
				glog.V(4).Infof("Processing %s", resource)
				c.ResourcesChan <- resource
				continue
			}

			glog.Warningf("Un-handled %v DeletedFinalStateUnknown for %s: %+v", delta.Type, deletedObject.Key, deletedObject.Obj)
			continue
		}

		glog.Warningf("Un-handled delta type %T (%s)", delta.Object, delta.Type)

	}
	return nil
}

// retry is a controller.RetryFunc that should return true if the given object and error
// should be retried after the provided number of times.
func (c *ExportController) retry(obj interface{}, err error, retries controller.Retry) bool {
	// let's retry a few times...
	return retries.Count < 5
}

// List is for the cache.ListerWatcher implementation
// List should return a list type object; the Items field will be extracted, and the
// ResourceVersion field will be used to start the watch in the right place.
func (c *ExportController) List(options kapi.ListOptions) (runtime.Object, error) {
	var err error
	options.LabelSelector, err = c.extendSelector(options.LabelSelector)
	if err != nil {
		return nil, err
	}

	glog.V(3).Infof("Running list func for %T with %+v", c.Kind, options)
	return c.ListFunc(options)
}

// Watch is for the cache.ListerWatcher implementation
// Watch should begin a watch at the specified version.
func (c *ExportController) Watch(options kapi.ListOptions) (watch.Interface, error) {
	var err error
	options.LabelSelector, err = c.extendSelector(options.LabelSelector)
	if err != nil {
		return nil, err
	}

	glog.V(3).Infof("Running watch func for %T with %+v", c.Kind, options)
	return c.WatchFunc(options)
}

// ListKeys implements the cache.KeyLister interface
// It is a function that returns the list of keys ("namespace/name" format)
// that we "know about" (to get a 2-way sync)
func (c *ExportController) ListKeys() []string {
	return c.KeyListFunc()
}

// GetByKey implements the cache.KeyGetter interface
// It is a function that returns the object that we "know about"
// for the given key ("namespace/name" format) - and a boolean if it exists
func (c *ExportController) GetByKey(key string) (interface{}, bool, error) {
	return c.KeyGetFunc(key)
}

// extendSelector extends the given labelSelector with the controller's
// requirements and user-provided labelSelector
func (c *ExportController) extendSelector(selector labels.Selector) (labels.Selector, error) {
	requirementFuncs := []func() (*labels.Requirement, error){}
	requirementFuncs = append(requirementFuncs, c.Requirements...)

	if len(c.LabelSelector) > 0 {
		requirements, err := labels.ParseToRequirements(c.LabelSelector)
		if err != nil {
			return nil, err
		}
		for i := range requirements {
			requirementFuncs = append(requirementFuncs, func() (*labels.Requirement, error) {
				return &requirements[i], nil
			})
		}
	}
	return extendSelector(selector, requirementFuncs...)
}
