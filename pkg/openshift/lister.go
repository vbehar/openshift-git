package openshift

import (
	"fmt"
	"reflect"

	"github.com/golang/glog"
	"github.com/openshift/origin/pkg/cmd/cli/cmd"

	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/cache"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/runtime"
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

	items, err := extractListItems(res)
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

// extractListItems extracts from the given object (which should be a *kapi.List)
// a slice of objects (contained in the Items field).
func extractListItems(object runtime.Object) ([]runtime.Object, error) {
	results := []runtime.Object{}

	objectValue := reflect.ValueOf(object)
	if !objectValue.IsValid() {
		return nil, fmt.Errorf("Can't get the value of %+v", object)
	}

	// object should be a *kapi.List
	if objectValue.Kind() == reflect.Ptr {
		objectValue = objectValue.Elem()
	}

	if objectValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("The list %+v should be a Struct, but it seems to be a %v", object, objectValue.Kind())
	}

	// the list should have an Items field
	itemsField := objectValue.FieldByName("Items")
	if !itemsField.IsValid() {
		return nil, fmt.Errorf("The struct %+v should have an 'Items' field", object)
	}

	if !itemsField.CanInterface() {
		return nil, fmt.Errorf("Can't interface the items field %v", itemsField)
	}
	items := itemsField.Interface()

	itemsValue := reflect.ValueOf(items)
	if !itemsValue.IsValid() {
		return nil, fmt.Errorf("Can't get the value of %+v", items)
	}

	// Items should be a slice
	if itemsValue.Kind() != reflect.Slice {
		return nil, fmt.Errorf("The Items %+v should be a Slice, but it seems to be a %v", items, itemsValue.Kind())
	}

	for i := 0; i < itemsValue.Len(); i++ {
		valueValue := itemsValue.Index(i)
		if !valueValue.IsValid() {
			return nil, fmt.Errorf("Can't get the value of elem %d in the Items %+v", i, items)
		}

		// The elems in the slice are not pointers,
		// but we need them to be pointers
		if valueValue.Kind() != reflect.Ptr && valueValue.CanAddr() {
			valueValue = valueValue.Addr()
		}
		if valueValue.Kind() != reflect.Ptr {
			return nil, fmt.Errorf("The elem %d in the Items %+v should be a Pointer, but it seems to be a %v", i, items, valueValue.Kind())
		}

		if !valueValue.CanInterface() {
			return nil, fmt.Errorf("Can't interface the elem %d %v", i, valueValue)
		}
		value := valueValue.Interface()

		if obj, ok := value.(runtime.Object); ok {
			results = append(results, obj)
		} else {
			return nil, fmt.Errorf("Value %+v is not a runtime.Object", value)
		}
	}

	return results, nil
}
