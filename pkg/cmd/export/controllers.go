package export

import (
	"github.com/vbehar/openshift-git/pkg/git"
	"github.com/vbehar/openshift-git/pkg/openshift"

	authorizationapi "github.com/openshift/origin/pkg/authorization/api"
	buildapi "github.com/openshift/origin/pkg/build/api"
	"github.com/openshift/origin/pkg/client"
	deployapi "github.com/openshift/origin/pkg/deploy/api"
	imageapi "github.com/openshift/origin/pkg/image/api"
	routeapi "github.com/openshift/origin/pkg/route/api"
	templateapi "github.com/openshift/origin/pkg/template/api"
	userapi "github.com/openshift/origin/pkg/user/api"

	kapi "k8s.io/kubernetes/pkg/api"
	kapiunversioned "k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/util/sets"
	"k8s.io/kubernetes/pkg/watch"

	"github.com/golang/glog"
)

func runBuildConfigsController(buildConfigs client.BuildConfigInterface,
	stopChan <-chan struct{}, resourcesChan chan<- openshift.Resource,
	repo *git.Repository, exportOptions *ExportOptions) {
	glog.V(1).Infof("Starting export controller for BuildConfigs !")
	(&openshift.ExportController{
		ResourcesChan: resourcesChan,
		LabelSelector: exportOptions.LabelSelector,
		ResyncPeriod:  exportOptions.ResyncPeriod,
		Kind:          &buildapi.BuildConfig{},
		KeyListFunc:   repo.KeyListFuncForKind(buildapi.BuildConfig{}),
		KeyGetFunc:    repo.KeyGetFuncForKindAndFormat(buildapi.BuildConfig{}, exportOptions.Format),
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return buildConfigs.List(options)
		},
		WatchFunc: func(options kapi.ListOptions) (watch.Interface, error) {
			return buildConfigs.Watch(options)
		},
	}).RunUntil(stopChan)
}

func runDeploymentConfigsController(deploymentConfigs client.DeploymentConfigInterface,
	stopChan <-chan struct{}, resourcesChan chan<- openshift.Resource,
	repo *git.Repository, exportOptions *ExportOptions) {
	glog.V(1).Infof("Starting export controller for DeploymentConfigs !")
	(&openshift.ExportController{
		ResourcesChan: resourcesChan,
		LabelSelector: exportOptions.LabelSelector,
		ResyncPeriod:  exportOptions.ResyncPeriod,
		Kind:          &deployapi.DeploymentConfig{},
		KeyListFunc:   repo.KeyListFuncForKind(deployapi.DeploymentConfig{}),
		KeyGetFunc:    repo.KeyGetFuncForKindAndFormat(deployapi.DeploymentConfig{}, exportOptions.Format),
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return deploymentConfigs.List(options)
		},
		WatchFunc: func(options kapi.ListOptions) (watch.Interface, error) {
			return deploymentConfigs.Watch(options)
		},
	}).RunUntil(stopChan)
}

func runImageStreamsController(imageStreams client.ImageStreamInterface,
	stopChan <-chan struct{}, resourcesChan chan<- openshift.Resource,
	repo *git.Repository, exportOptions *ExportOptions) {
	glog.V(1).Infof("Starting export controller for ImageStreams !")
	(&openshift.ExportController{
		ResourcesChan: resourcesChan,
		LabelSelector: exportOptions.LabelSelector,
		ResyncPeriod:  exportOptions.ResyncPeriod,
		Kind:          &imageapi.ImageStream{},
		KeyListFunc:   repo.KeyListFuncForKind(imageapi.ImageStream{}),
		KeyGetFunc:    repo.KeyGetFuncForKindAndFormat(imageapi.ImageStream{}, exportOptions.Format),
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return imageStreams.List(options)
		},
		WatchFunc: func(options kapi.ListOptions) (watch.Interface, error) {
			return imageStreams.Watch(options)
		},
	}).RunUntil(stopChan)
}

func runRoutesController(routes client.RouteInterface,
	stopChan <-chan struct{}, resourcesChan chan<- openshift.Resource,
	repo *git.Repository, exportOptions *ExportOptions) {
	glog.V(1).Infof("Starting export controller for Routes !")
	(&openshift.ExportController{
		ResourcesChan: resourcesChan,
		LabelSelector: exportOptions.LabelSelector,
		ResyncPeriod:  exportOptions.ResyncPeriod,
		Kind:          &routeapi.Route{},
		KeyListFunc:   repo.KeyListFuncForKind(routeapi.Route{}),
		KeyGetFunc:    repo.KeyGetFuncForKindAndFormat(routeapi.Route{}, exportOptions.Format),
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return routes.List(options)
		},
		WatchFunc: func(options kapi.ListOptions) (watch.Interface, error) {
			return routes.Watch(options)
		},
	}).RunUntil(stopChan)
}

func runTemplatesController(templates client.TemplateInterface,
	stopChan <-chan struct{}, resourcesChan chan<- openshift.Resource,
	repo *git.Repository, exportOptions *ExportOptions) {
	glog.V(1).Infof("Starting export controller for Templates !")
	(&openshift.ExportController{
		ResourcesChan: resourcesChan,
		LabelSelector: exportOptions.LabelSelector,
		ResyncPeriod:  exportOptions.ResyncPeriod,
		Kind:          &templateapi.Template{},
		KeyListFunc:   repo.KeyListFuncForKind(templateapi.Template{}),
		KeyGetFunc:    repo.KeyGetFuncForKindAndFormat(templateapi.Template{}, exportOptions.Format),
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return templates.List(options)
		},
		WatchFunc: func(options kapi.ListOptions) (watch.Interface, error) {
			return templates.Watch(options)
		},
	}).RunUntil(stopChan)
}

func runUsersController(users client.UserInterface,
	stopChan <-chan struct{}, resourcesChan chan<- openshift.Resource,
	repo *git.Repository, exportOptions *ExportOptions) {
	glog.V(1).Infof("Starting export controller for Users !")
	(&openshift.ExportController{
		ResourcesChan: resourcesChan,
		LabelSelector: exportOptions.LabelSelector,
		ResyncPeriod:  exportOptions.ResyncPeriod,
		Kind:          &userapi.User{},
		KeyListFunc:   repo.KeyListFuncForKind(userapi.User{}),
		KeyGetFunc:    repo.KeyGetFuncForKindAndFormat(userapi.User{}, exportOptions.Format),
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return users.List(options)
		},
		WatchFunc: func(options kapi.ListOptions) (watch.Interface, error) {
			return users.Watch(options)
		},
	}).RunUntil(stopChan)
}

func runGroupsController(groups client.GroupInterface,
	stopChan <-chan struct{}, resourcesChan chan<- openshift.Resource,
	repo *git.Repository, exportOptions *ExportOptions) {
	glog.V(1).Infof("Starting export controller for Groups !")
	(&openshift.ExportController{
		ResourcesChan: resourcesChan,
		LabelSelector: exportOptions.LabelSelector,
		ResyncPeriod:  exportOptions.ResyncPeriod,
		Kind:          &userapi.Group{},
		KeyListFunc:   repo.KeyListFuncForKind(userapi.Group{}),
		KeyGetFunc:    repo.KeyGetFuncForKindAndFormat(userapi.Group{}, exportOptions.Format),
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return groups.List(options)
		},
		WatchFunc: func(options kapi.ListOptions) (watch.Interface, error) {
			return groups.Watch(options)
		},
	}).RunUntil(stopChan)
}

func runPoliciesController(policies client.PolicyInterface,
	stopChan <-chan struct{}, resourcesChan chan<- openshift.Resource,
	repo *git.Repository, exportOptions *ExportOptions) {
	glog.V(1).Infof("Starting export controller for Policies !")
	(&openshift.ExportController{
		ResourcesChan: resourcesChan,
		LabelSelector: exportOptions.LabelSelector,
		ResyncPeriod:  exportOptions.ResyncPeriod,
		Kind:          &authorizationapi.Policy{},
		KeyListFunc:   repo.KeyListFuncForKind(authorizationapi.Policy{}),
		KeyGetFunc:    repo.KeyGetFuncForKindAndFormat(authorizationapi.Policy{}, exportOptions.Format),
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return policies.List(options)
		},
		WatchFunc: func(options kapi.ListOptions) (watch.Interface, error) {
			return policies.Watch(options)
		},
	}).RunUntil(stopChan)
}

func runPolicyBindingsController(policyBindings client.PolicyBindingInterface,
	stopChan <-chan struct{}, resourcesChan chan<- openshift.Resource,
	repo *git.Repository, exportOptions *ExportOptions) {
	glog.V(1).Infof("Starting export controller for PolicyBindings !")
	(&openshift.ExportController{
		ResourcesChan: resourcesChan,
		LabelSelector: exportOptions.LabelSelector,
		ResyncPeriod:  exportOptions.ResyncPeriod,
		Kind:          &authorizationapi.PolicyBinding{},
		KeyListFunc:   repo.KeyListFuncForKind(authorizationapi.PolicyBinding{}),
		KeyGetFunc:    repo.KeyGetFuncForKindAndFormat(authorizationapi.PolicyBinding{}, exportOptions.Format),
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return policyBindings.List(options)
		},
		WatchFunc: func(options kapi.ListOptions) (watch.Interface, error) {
			return policyBindings.Watch(options)
		},
	}).RunUntil(stopChan)
}

func runClusterPoliciesController(clusterPolicies client.ClusterPolicyInterface,
	stopChan <-chan struct{}, resourcesChan chan<- openshift.Resource,
	repo *git.Repository, exportOptions *ExportOptions) {
	glog.V(1).Infof("Starting export controller for ClusterPolicies !")
	(&openshift.ExportController{
		ResourcesChan: resourcesChan,
		LabelSelector: exportOptions.LabelSelector,
		ResyncPeriod:  exportOptions.ResyncPeriod,
		Kind:          &authorizationapi.ClusterPolicy{},
		KeyListFunc:   repo.KeyListFuncForKind(authorizationapi.ClusterPolicy{}),
		KeyGetFunc:    repo.KeyGetFuncForKindAndFormat(authorizationapi.ClusterPolicy{}, exportOptions.Format),
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return clusterPolicies.List(options)
		},
		WatchFunc: func(options kapi.ListOptions) (watch.Interface, error) {
			return clusterPolicies.Watch(options)
		},
	}).RunUntil(stopChan)
}

func runClusterPolicyBindingsController(clusterPolicyBindings client.ClusterPolicyBindingInterface,
	stopChan <-chan struct{}, resourcesChan chan<- openshift.Resource,
	repo *git.Repository, exportOptions *ExportOptions) {
	glog.V(1).Infof("Starting export controller for ClusterPolicyBindings !")
	(&openshift.ExportController{
		ResourcesChan: resourcesChan,
		LabelSelector: exportOptions.LabelSelector,
		ResyncPeriod:  exportOptions.ResyncPeriod,
		Kind:          &authorizationapi.ClusterPolicyBinding{},
		KeyListFunc:   repo.KeyListFuncForKind(authorizationapi.ClusterPolicyBinding{}),
		KeyGetFunc:    repo.KeyGetFuncForKindAndFormat(authorizationapi.ClusterPolicyBinding{}, exportOptions.Format),
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return clusterPolicyBindings.List(options)
		},
		WatchFunc: func(options kapi.ListOptions) (watch.Interface, error) {
			return clusterPolicyBindings.Watch(options)
		},
	}).RunUntil(stopChan)
}

func runReplicationControllersController(replicationControllers unversioned.ReplicationControllerInterface,
	stopChan <-chan struct{}, resourcesChan chan<- openshift.Resource,
	repo *git.Repository, exportOptions *ExportOptions) {
	glog.V(1).Infof("Starting export controller for ReplicationControllers !")
	(&openshift.ExportController{
		ResourcesChan: resourcesChan,
		LabelSelector: exportOptions.LabelSelector,
		ResyncPeriod:  exportOptions.ResyncPeriod,
		Kind:          &kapi.ReplicationController{},
		KeyListFunc:   repo.KeyListFuncForKind(kapi.ReplicationController{}),
		KeyGetFunc:    repo.KeyGetFuncForKindAndFormat(kapi.ReplicationController{}, exportOptions.Format),
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return replicationControllers.List(options)
		},
		WatchFunc: func(options kapi.ListOptions) (watch.Interface, error) {
			return replicationControllers.Watch(options)
		},
		Requirements: []func() (*labels.Requirement, error){
			func() (*labels.Requirement, error) {
				return labels.NewRequirement(deployapi.DeploymentConfigAnnotation, labels.DoesNotExistOperator, sets.NewString())
			},
		},
	}).RunUntil(stopChan)
}

func runPodsController(pods unversioned.PodInterface,
	stopChan <-chan struct{}, resourcesChan chan<- openshift.Resource,
	repo *git.Repository, exportOptions *ExportOptions) {
	glog.V(1).Infof("Starting export controller for Pods !")
	(&openshift.ExportController{
		ResourcesChan: resourcesChan,
		LabelSelector: exportOptions.LabelSelector,
		ResyncPeriod:  exportOptions.ResyncPeriod,
		Kind:          &kapi.Pod{},
		KeyListFunc:   repo.KeyListFuncForKind(kapi.Pod{}),
		KeyGetFunc:    repo.KeyGetFuncForKindAndFormat(kapi.Pod{}, exportOptions.Format),
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return pods.List(options)
		},
		WatchFunc: func(options kapi.ListOptions) (watch.Interface, error) {
			return pods.Watch(options)
		},
		Requirements: []func() (*labels.Requirement, error){
			func() (*labels.Requirement, error) {
				return labels.NewRequirement(buildapi.BuildLabel, labels.DoesNotExistOperator, sets.NewString())
			},
			func() (*labels.Requirement, error) {
				return labels.NewRequirement(deployapi.DeployerPodForDeploymentLabel, labels.DoesNotExistOperator, sets.NewString())
			},
			func() (*labels.Requirement, error) {
				return labels.NewRequirement(deployapi.DeploymentConfigLabel, labels.DoesNotExistOperator, sets.NewString())
			},
		},
	}).RunUntil(stopChan)
}

func runEndpointsController(endpoints unversioned.EndpointsInterface,
	stopChan <-chan struct{}, resourcesChan chan<- openshift.Resource,
	repo *git.Repository, exportOptions *ExportOptions) {
	glog.V(1).Infof("Starting export controller for Endpoints !")
	(&openshift.ExportController{
		ResourcesChan: resourcesChan,
		LabelSelector: exportOptions.LabelSelector,
		ResyncPeriod:  exportOptions.ResyncPeriod,
		Kind:          &kapi.Endpoints{},
		KeyListFunc:   repo.KeyListFuncForKind(kapi.Endpoints{}),
		KeyGetFunc:    repo.KeyGetFuncForKindAndFormat(kapi.Endpoints{}, exportOptions.Format),
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return endpoints.List(options)
		},
		WatchFunc: func(options kapi.ListOptions) (watch.Interface, error) {
			return endpoints.Watch(options)
		},
	}).RunUntil(stopChan)
}

func runServicesController(services unversioned.ServiceInterface,
	stopChan <-chan struct{}, resourcesChan chan<- openshift.Resource,
	repo *git.Repository, exportOptions *ExportOptions) {
	glog.V(1).Infof("Starting export controller for Services !")
	(&openshift.ExportController{
		ResourcesChan: resourcesChan,
		LabelSelector: exportOptions.LabelSelector,
		ResyncPeriod:  exportOptions.ResyncPeriod,
		Kind:          &kapi.Service{},
		KeyListFunc:   repo.KeyListFuncForKind(kapi.Service{}),
		KeyGetFunc:    repo.KeyGetFuncForKindAndFormat(kapi.Service{}, exportOptions.Format),
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return services.List(options)
		},
		WatchFunc: func(options kapi.ListOptions) (watch.Interface, error) {
			return services.Watch(options)
		},
	}).RunUntil(stopChan)
}

func runServiceAccountsController(serviceAccounts unversioned.ServiceAccountsInterface,
	stopChan <-chan struct{}, resourcesChan chan<- openshift.Resource,
	repo *git.Repository, exportOptions *ExportOptions) {
	glog.V(1).Infof("Starting export controller for ServiceAccounts !")
	(&openshift.ExportController{
		ResourcesChan: resourcesChan,
		LabelSelector: exportOptions.LabelSelector,
		ResyncPeriod:  exportOptions.ResyncPeriod,
		Kind:          &kapi.ServiceAccount{},
		KeyListFunc:   repo.KeyListFuncForKind(kapi.ServiceAccount{}),
		KeyGetFunc:    repo.KeyGetFuncForKindAndFormat(kapi.ServiceAccount{}, exportOptions.Format),
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return serviceAccounts.List(options)
		},
		WatchFunc: func(options kapi.ListOptions) (watch.Interface, error) {
			return serviceAccounts.Watch(options)
		},
	}).RunUntil(stopChan)
}

func runSecretsController(secrets unversioned.SecretsInterface,
	stopChan <-chan struct{}, resourcesChan chan<- openshift.Resource,
	repo *git.Repository, exportOptions *ExportOptions) {
	glog.V(1).Infof("Starting export controller for Secrets !")
	(&openshift.ExportController{
		ResourcesChan: resourcesChan,
		LabelSelector: exportOptions.LabelSelector,
		ResyncPeriod:  exportOptions.ResyncPeriod,
		Kind:          &kapi.Secret{},
		KeyListFunc:   repo.KeyListFuncForKind(kapi.Secret{}),
		KeyGetFunc:    repo.KeyGetFuncForKindAndFormat(kapi.Secret{}, exportOptions.Format),
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return secrets.List(options)
		},
		WatchFunc: func(options kapi.ListOptions) (watch.Interface, error) {
			return secrets.Watch(options)
		},
	}).RunUntil(stopChan)
}

func runLimitRangesController(limitRanges unversioned.LimitRangeInterface,
	stopChan <-chan struct{}, resourcesChan chan<- openshift.Resource,
	repo *git.Repository, exportOptions *ExportOptions) {
	glog.V(1).Infof("Starting export controller for LimitRanges !")
	(&openshift.ExportController{
		ResourcesChan: resourcesChan,
		LabelSelector: exportOptions.LabelSelector,
		ResyncPeriod:  exportOptions.ResyncPeriod,
		Kind:          &kapi.LimitRange{},
		KeyListFunc:   repo.KeyListFuncForKind(kapi.LimitRange{}),
		KeyGetFunc:    repo.KeyGetFuncForKindAndFormat(kapi.LimitRange{}, exportOptions.Format),
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return limitRanges.List(options)
		},
		WatchFunc: func(options kapi.ListOptions) (watch.Interface, error) {
			return limitRanges.Watch(options)
		},
	}).RunUntil(stopChan)
}

func runResourceQuotasController(resourceQuotas unversioned.ResourceQuotaInterface,
	stopChan <-chan struct{}, resourcesChan chan<- openshift.Resource,
	repo *git.Repository, exportOptions *ExportOptions) {
	glog.V(1).Infof("Starting export controller for ResourceQuotas !")
	(&openshift.ExportController{
		ResourcesChan: resourcesChan,
		LabelSelector: exportOptions.LabelSelector,
		ResyncPeriod:  exportOptions.ResyncPeriod,
		Kind:          &kapi.ResourceQuota{},
		KeyListFunc:   repo.KeyListFuncForKind(kapi.ResourceQuota{}),
		KeyGetFunc:    repo.KeyGetFuncForKindAndFormat(kapi.ResourceQuota{}, exportOptions.Format),
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return resourceQuotas.List(options)
		},
		WatchFunc: func(options kapi.ListOptions) (watch.Interface, error) {
			return resourceQuotas.Watch(options)
		},
	}).RunUntil(stopChan)
}

func runPersistentVolumeClaimsController(persistentVolumeClaims unversioned.PersistentVolumeClaimInterface,
	stopChan <-chan struct{}, resourcesChan chan<- openshift.Resource,
	repo *git.Repository, exportOptions *ExportOptions) {
	glog.V(1).Infof("Starting export controller for PersistentVolumeClaims !")
	(&openshift.ExportController{
		ResourcesChan: resourcesChan,
		LabelSelector: exportOptions.LabelSelector,
		ResyncPeriod:  exportOptions.ResyncPeriod,
		Kind:          &kapi.PersistentVolumeClaim{},
		KeyListFunc:   repo.KeyListFuncForKind(kapi.PersistentVolumeClaim{}),
		KeyGetFunc:    repo.KeyGetFuncForKindAndFormat(kapi.PersistentVolumeClaim{}, exportOptions.Format),
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return persistentVolumeClaims.List(options)
		},
		WatchFunc: func(options kapi.ListOptions) (watch.Interface, error) {
			return persistentVolumeClaims.Watch(options)
		},
	}).RunUntil(stopChan)
}

func runPersistentVolumesController(persistentVolumes unversioned.PersistentVolumeInterface,
	stopChan <-chan struct{}, resourcesChan chan<- openshift.Resource,
	repo *git.Repository, exportOptions *ExportOptions) {
	glog.V(1).Infof("Starting export controller for PersistentVolumes !")
	(&openshift.ExportController{
		ResourcesChan: resourcesChan,
		LabelSelector: exportOptions.LabelSelector,
		ResyncPeriod:  exportOptions.ResyncPeriod,
		Kind:          &kapi.PersistentVolume{},
		KeyListFunc:   repo.KeyListFuncForKind(kapi.PersistentVolume{}),
		KeyGetFunc:    repo.KeyGetFuncForKindAndFormat(kapi.PersistentVolume{}, exportOptions.Format),
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return persistentVolumes.List(options)
		},
		WatchFunc: func(options kapi.ListOptions) (watch.Interface, error) {
			return persistentVolumes.Watch(options)
		},
	}).RunUntil(stopChan)
}

func runSecurityContextConstraintsController(securityContextConstraints unversioned.SecurityContextConstraintInterface,
	stopChan <-chan struct{}, resourcesChan chan<- openshift.Resource,
	repo *git.Repository, exportOptions *ExportOptions) {
	glog.V(1).Infof("Starting export controller for SecurityContextConstraints !")
	(&openshift.ExportController{
		ResourcesChan: resourcesChan,
		LabelSelector: exportOptions.LabelSelector,
		ResyncPeriod:  exportOptions.ResyncPeriod,
		Kind:          &kapi.SecurityContextConstraints{},
		KeyListFunc:   repo.KeyListFuncForKind(kapi.SecurityContextConstraints{}),
		KeyGetFunc:    repo.KeyGetFuncForKindAndFormat(kapi.SecurityContextConstraints{}, exportOptions.Format),
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return securityContextConstraints.List(options)
		},
		WatchFunc: func(options kapi.ListOptions) (watch.Interface, error) {
			return securityContextConstraints.Watch(options)
		},
	}).RunUntil(stopChan)
}

func runNamespacesController(namespaces unversioned.NamespaceInterface,
	stopChan <-chan struct{}, resourcesChan chan<- openshift.Resource,
	repo *git.Repository, exportOptions *ExportOptions) {
	glog.V(1).Infof("Starting export controller for Namespaces !")
	(&openshift.ExportController{
		ResourcesChan: resourcesChan,
		LabelSelector: exportOptions.LabelSelector,
		ResyncPeriod:  exportOptions.ResyncPeriod,
		Kind:          &kapi.Namespace{},
		KeyListFunc:   repo.KeyListFuncForKind(kapi.Namespace{}),
		KeyGetFunc:    repo.KeyGetFuncForKindAndFormat(kapi.Namespace{}, exportOptions.Format),
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return namespaces.List(options)
		},
		WatchFunc: func(options kapi.ListOptions) (watch.Interface, error) {
			return namespaces.Watch(options)
		},
	}).RunUntil(stopChan)
}

func runNamespaceController(namespaces unversioned.NamespaceInterface, namespace string,
	stopChan <-chan struct{}, resourcesChan chan<- openshift.Resource,
	repo *git.Repository, exportOptions *ExportOptions) {
	glog.V(1).Infof("Starting export controller for Namespace %s !", namespace)
	(&openshift.ExportController{
		ResourcesChan: resourcesChan,
		ResyncPeriod:  exportOptions.ResyncPeriod,
		Kind:          &kapi.Namespace{},
		KeyListFunc:   repo.KeyListFuncForKind(kapi.Namespace{}),
		KeyGetFunc:    repo.KeyGetFuncForKindAndFormat(kapi.Namespace{}, exportOptions.Format),
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			ns, err := namespaces.Get(namespace)
			if err != nil {
				return nil, err
			}
			list := &kapi.NamespaceList{
				TypeMeta: kapiunversioned.TypeMeta{
					Kind:       "List",
					APIVersion: ns.APIVersion,
				},
				Items: []kapi.Namespace{
					*ns,
				},
			}
			return list, nil
		},
		WatchFunc: func(options kapi.ListOptions) (watch.Interface, error) {
			// can't watch a single specific namespace, so let's watch nothing for the moment
			return watch.NewFake(), nil
		},
	}).RunUntil(stopChan)
}
