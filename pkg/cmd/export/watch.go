package export

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/vbehar/openshift-git/pkg/git"
	"github.com/vbehar/openshift-git/pkg/openshift"

	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/meta"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/kubectl"
	"k8s.io/kubernetes/pkg/kubectl/resource"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/watch"

	"github.com/golang/glog"
)

// runWatch run the export controllers for the given resources
func runWatch(resources string, repo *git.Repository) error {
	saveWaiter := &sync.WaitGroup{}
	stopChan := make(chan struct{})
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

		if mapping.Scope.Name() == meta.RESTScopeNameRoot && !exportOptions.AllNamespaces {
			switch gvk.Kind {
			case "Namespace", "Project":
				if err := runControllerForNamespace(gvk, namespace, mapper, restClient, stopChan, resourcesChan, repo, exportOptions); err != nil {
					return err
				}
			default:
				glog.Warningf("Ignoring root kind %s because you asked for a specific namespace", gvk)
			}
		} else {
			if err := runController(gvk, namespace, mapper, restClient, stopChan, resourcesChan, repo, exportOptions); err != nil {
				return err
			}
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	select {
	case <-c:
		glog.Infof("Interrupted by user (or killed) !")
		close(stopChan)
		time.Sleep(1 * time.Second)
		close(resourcesChan)
		saveWaiter.Wait()
	}

	return nil
}

// runController starts an export controller (in a new goroutine) for the given kind,
// in the given namespace
func runController(gvk unversioned.GroupVersionKind,
	namespace string,
	mapper meta.RESTMapper, restClient resource.RESTClient,
	stopChan <-chan struct{}, resourcesChan chan<- openshift.Resource,
	repo *git.Repository, exportOptions *ExportOptions) error {

	if !kapi.Scheme.Recognizes(gvk) {
		return fmt.Errorf("GVK %s not recognizes", gvk)
	}

	obj, err := kapi.Scheme.New(gvk)
	if err != nil {
		return err
	}

	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return err
	}

	helper := resource.NewHelper(restClient, mapping)

	var requirements []func() (*labels.Requirement, error)
	if exportOptions.UseDefaultSelector {
		requirements = DefaultRequirementsFor(gvk)
	}

	glog.V(1).Infof("Starting export controller for %s", gvk.Kind)
	(&openshift.ExportController{
		ResourcesChan: resourcesChan,
		LabelSelector: exportOptions.LabelSelector,
		ResyncPeriod:  exportOptions.ResyncPeriod,
		Kind:          obj,
		KeyListFunc:   repo.KeyListFuncForKind(gvk.Kind),
		KeyGetFunc:    repo.KeyGetFuncForKindAndFormat(gvk.Kind, exportOptions.Format),
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return helper.List(namespace, gvk.Version, options.LabelSelector, false)
		},
		WatchFunc: func(options kapi.ListOptions) (watch.Interface, error) {
			return helper.Watch(namespace, options.ResourceVersion, gvk.Version, options.LabelSelector)
		},
		Requirements: requirements,
	}).RunUntil(stopChan)

	return nil
}

// runControllerForNamespace starts an export controller (in a new goroutine)
// that can be used to export a single namespace/project.
func runControllerForNamespace(gvk unversioned.GroupVersionKind,
	namespace string,
	mapper meta.RESTMapper, restClient resource.RESTClient,
	stopChan <-chan struct{}, resourcesChan chan<- openshift.Resource,
	repo *git.Repository, exportOptions *ExportOptions) error {

	gvkList := gvk.GroupVersion().WithKind(gvk.Kind + "List")

	if !kapi.Scheme.Recognizes(gvk) {
		return fmt.Errorf("GVK %s not recognizes", gvk)
	}

	obj, err := kapi.Scheme.New(gvk)
	if err != nil {
		return err
	}

	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return err
	}

	helper := resource.NewHelper(restClient, mapping)

	var requirements []func() (*labels.Requirement, error)
	if exportOptions.UseDefaultSelector {
		requirements = DefaultRequirementsFor(gvk)
	}

	glog.V(1).Infof("Starting export controller for %s %s...", gvk.Kind, namespace)
	(&openshift.ExportController{
		ResourcesChan: resourcesChan,
		LabelSelector: exportOptions.LabelSelector,
		ResyncPeriod:  exportOptions.ResyncPeriod,
		Kind:          obj,
		KeyListFunc:   repo.KeyListFuncForKind(gvk.Kind),
		KeyGetFunc:    repo.KeyGetFuncForKindAndFormat(gvk.Kind, exportOptions.Format),
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
		WatchFunc: func(options kapi.ListOptions) (watch.Interface, error) {
			// can't watch a single specific namespace, so let's watch nothing for the moment
			return watch.NewFake(), nil
		},
		Requirements: requirements,
	}).RunUntil(stopChan)

	return nil
}
