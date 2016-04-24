package export

import (
	"fmt"
	"strings"
	"sync"

	"github.com/vbehar/openshift-git/pkg/git"
	"github.com/vbehar/openshift-git/pkg/openshift"

	"github.com/openshift/origin/pkg/util/parallel"

	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/meta"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/kubectl"
	"k8s.io/kubernetes/pkg/kubectl/resource"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/runtime"

	"github.com/golang/glog"
)

// runList run the "list" operations in parallel to export the given resources to the given repository
func runList(resources string, repo *git.Repository) error {
	saveWaiter := &sync.WaitGroup{}
	resourcesChan := make(chan openshift.Resource, 10)

	namespace, _, err := openshift.Factory.DefaultNamespace()
	if err != nil {
		return err
	}
	if exportOptions.AllNamespaces {
		namespace = kapi.NamespaceAll
	}

	mapper, _ := openshift.Factory.Object()
	oclient, kclient, err := openshift.Factory.Clients()
	if err != nil {
		return err
	}

	kinds, err := openshift.KindsFor(mapper, resource.SplitResourceArgument(resources))
	if err != nil {
		return err
	}
	if len(kinds) == 0 {
		return fmt.Errorf("No valid kinds for '%s'", resources)
	}

	if exportOptions.AllNamespaces {
		glog.Infof("Running export for kinds %v for all namespaces", kinds)
	} else {
		glog.Infof("Running export for kinds %v for namespace %s", kinds, namespace)
	}

	printer, _, err := kubectl.GetPrinter(exportOptions.Format, "")
	if err != nil {
		return err
	}

	saveWaiter.Add(1)
	go func() {
		defer saveWaiter.Done()
		saveResources(repo, resourcesChan, mapper, printer)
	}()

	knownTypes := kapi.Scheme.KnownTypes(kapi.Unversioned)

	listers := []func() error{}
	for _, gvk := range kinds {
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			return err
		}

		var restClient resource.RESTClient
		restClient = kclient
		if kindType, found := knownTypes[gvk.Kind]; found {
			if strings.Contains(kindType.PkgPath(), "openshift") {
				restClient = oclient
			}
		}

		var lister func() error
		if mapping.Scope.Name() == meta.RESTScopeNameRoot && !exportOptions.AllNamespaces {
			switch gvk.Kind {
			case "Namespace", "Project":
				lister = listerForNamespace(gvk, namespace, mapper, restClient, resourcesChan, exportOptions.LabelSelector)
			default:
				glog.Warningf("Ignoring root kind %s because you asked for a specific namespace", gvk)
			}
		} else {
			lister = listerFor(gvk, namespace, mapper, restClient, resourcesChan, exportOptions.LabelSelector)
		}

		if lister != nil {
			listers = append(listers, lister)
		}
	}

	if errs := parallel.Run(listers...); len(errs) > 0 {
		return fmt.Errorf("Got %d errors: %+v", len(errs), errs)
	}

	close(resourcesChan)
	saveWaiter.Wait()

	return nil
}

// listerFor returns a "lister" func that can be used to list objects of the given kind,
// in the given namespace
func listerFor(gvk unversioned.GroupVersionKind,
	namespace string,
	mapper meta.RESTMapper, restClient resource.RESTClient,
	resourcesChan chan<- openshift.Resource, labelSelector string) func() error {

	if !kapi.Scheme.Recognizes(gvk) {
		return func() error { return fmt.Errorf("GVK %s not recognizes", gvk) }
	}

	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return func() error { return err }
	}

	helper := resource.NewHelper(restClient, mapping)

	var requirements []func() (*labels.Requirement, error)
	if exportOptions.UseDefaultSelector {
		requirements = DefaultRequirementsFor(gvk)
	}

	glog.V(1).Infof("Listing %s...", gvk.Kind)
	return (&openshift.ExportLister{
		ResourcesChan: resourcesChan,
		LabelSelector: labelSelector,
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return helper.List(namespace, gvk.Version, options.LabelSelector, false)
		},
		Requirements: requirements,
	}).List
}

// listerForNamespace returns a "lister" func that can be used to list a single namespace/project.
func listerForNamespace(gvk unversioned.GroupVersionKind,
	namespace string,
	mapper meta.RESTMapper, restClient resource.RESTClient,
	resourcesChan chan<- openshift.Resource, labelSelector string) func() error {

	gvkList := gvk.GroupVersion().WithKind(gvk.Kind + "List")

	if !kapi.Scheme.Recognizes(gvk) {
		return func() error { return fmt.Errorf("GVK %s not recognizes", gvk) }
	}

	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return func() error { return err }
	}

	helper := resource.NewHelper(restClient, mapping)

	var requirements []func() (*labels.Requirement, error)
	if exportOptions.UseDefaultSelector {
		requirements = DefaultRequirementsFor(gvk)
	}

	glog.V(1).Infof("Getting %s %s...", gvk.Kind, namespace)
	return (&openshift.ExportLister{
		ResourcesChan: resourcesChan,
		LabelSelector: labelSelector,
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			obj, err := helper.Get(namespace, namespace, false)
			if err != nil {
				return nil, err
			}

			newObj, err := kapi.Scheme.ConvertToVersion(obj, gvk.Version)
			if err != nil {
				return nil, err
			}

			listObject, err := kapi.Scheme.New(gvkList)
			if err != nil {
				return nil, err
			}

			if err := meta.SetList(listObject, []runtime.Object{newObj}); err != nil {
				return nil, err
			}

			return listObject, nil
		},
		Requirements: requirements,
	}).List
}
