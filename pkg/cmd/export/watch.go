package export

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/vbehar/openshift-git/pkg/git"
	"github.com/vbehar/openshift-git/pkg/openshift"

	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/kubectl"

	"github.com/golang/glog"
)

func runWatch(repo *git.Repository) {
	saveWaiter := &sync.WaitGroup{}
	stopChan := make(chan struct{})
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

	if exportOptions.AllNamespaces {
		runNamespacesController(kclient.Namespaces(), stopChan, resourcesChan, repo, exportOptions)
		runPersistentVolumesController(kclient.PersistentVolumes(), stopChan, resourcesChan, repo, exportOptions)
		runSecurityContextConstraintsController(kclient.SecurityContextConstraints(), stopChan, resourcesChan, repo, exportOptions)
		runClusterPoliciesController(oclient.ClusterPolicies(), stopChan, resourcesChan, repo, exportOptions)
		runClusterPolicyBindingsController(oclient.ClusterPolicyBindings(), stopChan, resourcesChan, repo, exportOptions)
		runUsersController(oclient.Users(), stopChan, resourcesChan, repo, exportOptions)
		runGroupsController(oclient.Groups(), stopChan, resourcesChan, repo, exportOptions)
	} else {
		runNamespaceController(kclient.Namespaces(), namespace, stopChan, resourcesChan, repo, exportOptions)
	}

	runBuildConfigsController(oclient.BuildConfigs(namespace), stopChan, resourcesChan, repo, exportOptions)
	runDeploymentConfigsController(oclient.DeploymentConfigs(namespace), stopChan, resourcesChan, repo, exportOptions)
	runReplicationControllersController(kclient.ReplicationControllers(namespace), stopChan, resourcesChan, repo, exportOptions)
	runPodsController(kclient.Pods(namespace), stopChan, resourcesChan, repo, exportOptions)
	runImageStreamsController(oclient.ImageStreams(namespace), stopChan, resourcesChan, repo, exportOptions)
	runServicesController(kclient.Services(namespace), stopChan, resourcesChan, repo, exportOptions)
	runRoutesController(oclient.Routes(namespace), stopChan, resourcesChan, repo, exportOptions)
	runTemplatesController(oclient.Templates(namespace), stopChan, resourcesChan, repo, exportOptions)
	runSecretsController(kclient.Secrets(namespace), stopChan, resourcesChan, repo, exportOptions)
	runLimitRangesController(kclient.LimitRanges(namespace), stopChan, resourcesChan, repo, exportOptions)
	runResourceQuotasController(kclient.ResourceQuotas(namespace), stopChan, resourcesChan, repo, exportOptions)
	runPersistentVolumeClaimsController(kclient.PersistentVolumeClaims(namespace), stopChan, resourcesChan, repo, exportOptions)
	runPoliciesController(oclient.Policies(namespace), stopChan, resourcesChan, repo, exportOptions)
	runPolicyBindingsController(oclient.PolicyBindings(namespace), stopChan, resourcesChan, repo, exportOptions)
	runServiceAccountsController(kclient.ServiceAccounts(namespace), stopChan, resourcesChan, repo, exportOptions)

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
}
