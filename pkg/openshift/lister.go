package openshift

import (
	"github.com/openshift/origin/pkg/cmd/cli/cmd"

	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/meta"
	"k8s.io/kubernetes/pkg/client/cache"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/runtime"

	"github.com/golang/glog"
)

// ExportLister represents a lister that can list some resources,
// and push an "exported" view of these resources to a channel.
type ExportLister struct {
	// ResourcesChan is the channel over which resources will be send
	ResourcesChan chan<- Resource

	// ListFunc is the function used to list the resources for the provided kind
	ListFunc func(options kapi.ListOptions) (runtime.Object, error)

	// Requirements is a slice of functions that can create requirements
	// for customizing the labelSelector for the provided kind
	Requirements []func() (*labels.Requirement, error)

	// LabelSelector is a user-provided labelSelector as a string
	// used to restrict the resources for the provided kind
	LabelSelector string
}

// List lists the resources and push them to the channel
func (l *ExportLister) List() error {
	var err error
	exporter := cmd.NewExporter()

	options := &kapi.ListOptions{}
	options.LabelSelector, err = l.extendSelector(labels.Everything())
	if err != nil {
		return err
	}

	res, err := l.ListFunc(*options)
	if err != nil {
		return err
	}

	items, err := meta.ExtractList(res)
	if err != nil {
		return err
	}
	glog.V(2).Infof("Found %d items for %T", len(items), res)

	for _, obj := range items {
		glog.V(5).Infof("Handling %T", obj)

		// get the reference before exporting
		// (after it will be too late to get a reference)
		var ref *kapi.ObjectReference
		ref, err = kapi.GetReference(obj)
		if err != nil {
			return err
		}

		if err := exporter.Export(obj, false); err != nil {
			if err == cmd.ErrExportOmit {
				// let's just ignore this object that can't be exported
				glog.V(6).Infof("Ignoring %s %s/%s: %v", ref.Kind, ref.Namespace, ref.Name, err)
				continue
			}
			return err
		}

		r := Resource{
			ObjectReference: ref,
			Object:          obj,
			Exists:          true,
			Status:          string(cache.Sync),
		}

		glog.V(4).Infof("Processing %s", r)
		l.ResourcesChan <- r
	}

	return nil
}

// extendSelector extends the given labelSelector with the lister's
// requirements and user-provided labelSelector
func (l *ExportLister) extendSelector(selector labels.Selector) (labels.Selector, error) {
	requirementFuncs := []func() (*labels.Requirement, error){}
	requirementFuncs = append(requirementFuncs, l.Requirements...)

	if len(l.LabelSelector) > 0 {
		requirements, err := labels.ParseToRequirements(l.LabelSelector)
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
