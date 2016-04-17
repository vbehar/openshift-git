package export

import (
	"sync"

	"github.com/vbehar/openshift-git/pkg/git"
	"github.com/vbehar/openshift-git/pkg/openshift"

	"github.com/openshift/origin/pkg/util/parallel"

	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/kubectl"

	"github.com/golang/glog"
)

func runList(repo *git.Repository) {
	saveWaiter := &sync.WaitGroup{}
	resourcesChan := make(chan openshift.Resource, 10)

	namespace, _, err := openshift.Factory.DefaultNamespace()
	if err != nil {
		glog.Fatalf("Failed to get namespace: %v", err)
	}
	if exportOptions.AllNamespaces {
		namespace = kapi.NamespaceAll
	}

	mapper, _ := openshift.Factory.Object()
	oclient, kclient, err := openshift.Factory.Clients()
	if err != nil {
		glog.Fatalf("Failed to get openshift client from factory: %v", err)
	}

	printer, _, err := kubectl.GetPrinter(exportOptions.Format, "")
	if err != nil {
		glog.Fatalf("Failed to get printer for format %s: %v", exportOptions.Format, err)
	}

	saveWaiter.Add(1)
	go func() {
		defer saveWaiter.Done()
		saveResources(repo, resourcesChan, mapper, printer)
	}()

	allErrors := []error{}
	if exportOptions.AllNamespaces {
		errs := parallel.Run(
			namespacesLister(kclient.Namespaces(), resourcesChan, exportOptions.LabelSelector),
			persistentVolumesLister(kclient.PersistentVolumes(), resourcesChan, exportOptions.LabelSelector),
			securityContextConstraintsLister(kclient.SecurityContextConstraints(), resourcesChan, exportOptions.LabelSelector),
			clusterPoliciesLister(oclient.ClusterPolicies(), resourcesChan, exportOptions.LabelSelector),
			clusterPolicyBindingsLister(oclient.ClusterPolicyBindings(), resourcesChan, exportOptions.LabelSelector),
			usersLister(oclient.Users(), resourcesChan, exportOptions.LabelSelector),
			groupsLister(oclient.Groups(), resourcesChan, exportOptions.LabelSelector),
		)
		allErrors = append(allErrors, errs...)
	} else {
		errs := parallel.Run(
			namespaceLister(kclient.Namespaces(), namespace, resourcesChan, exportOptions.LabelSelector),
		)
		allErrors = append(allErrors, errs...)
	}

	errs := parallel.Run(
		buildConfigsLister(oclient.BuildConfigs(namespace), resourcesChan, exportOptions.LabelSelector),
		deploymentConfigsLister(oclient.DeploymentConfigs(namespace), resourcesChan, exportOptions.LabelSelector),
		replicationControllersLister(kclient.ReplicationControllers(namespace), resourcesChan, exportOptions.LabelSelector),
		podsLister(kclient.Pods(namespace), resourcesChan, exportOptions.LabelSelector),
		imageStreamsLister(oclient.ImageStreams(namespace), resourcesChan, exportOptions.LabelSelector),
		servicesLister(kclient.Services(namespace), resourcesChan, exportOptions.LabelSelector),
		routesLister(oclient.Routes(namespace), resourcesChan, exportOptions.LabelSelector),
		templatesLister(oclient.Templates(namespace), resourcesChan, exportOptions.LabelSelector),
		secretsLister(kclient.Secrets(namespace), resourcesChan, exportOptions.LabelSelector),
		limitRangesLister(kclient.LimitRanges(namespace), resourcesChan, exportOptions.LabelSelector),
		resourceQuotasLister(kclient.ResourceQuotas(namespace), resourcesChan, exportOptions.LabelSelector),
		persistentVolumeClaimsLister(kclient.PersistentVolumeClaims(namespace), resourcesChan, exportOptions.LabelSelector),
		policiesLister(oclient.Policies(namespace), resourcesChan, exportOptions.LabelSelector),
		policyBindingsLister(oclient.PolicyBindings(namespace), resourcesChan, exportOptions.LabelSelector),
		serviceAccountsLister(kclient.ServiceAccounts(namespace), resourcesChan, exportOptions.LabelSelector),
	)
	allErrors = append(allErrors, errs...)

	if len(allErrors) > 0 {
		glog.Fatalf("Got %d errors: %+v", len(allErrors), allErrors)
	}

	close(resourcesChan)
	saveWaiter.Wait()
}
