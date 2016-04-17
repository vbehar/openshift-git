package export

import (
	"github.com/golang/glog"
	"github.com/vbehar/openshift-git/pkg/openshift"

	buildapi "github.com/openshift/origin/pkg/build/api"
	"github.com/openshift/origin/pkg/client"
	deployapi "github.com/openshift/origin/pkg/deploy/api"

	kapi "k8s.io/kubernetes/pkg/api"
	kapiunversioned "k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/util/sets"
)

func buildConfigsLister(buildConfigs client.BuildConfigInterface, resourcesChan chan<- openshift.Resource, labelSelector string) func() error {
	glog.V(1).Infof("Listing BuildConfigs...")
	return (&openshift.ExportLister{
		ResourcesChan: resourcesChan,
		LabelSelector: labelSelector,
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return buildConfigs.List(options)
		},
	}).List
}

func deploymentConfigsLister(deploymentConfigs client.DeploymentConfigInterface, resourcesChan chan<- openshift.Resource, labelSelector string) func() error {
	glog.V(1).Infof("Listing DeploymentConfigs...")
	return (&openshift.ExportLister{
		ResourcesChan: resourcesChan,
		LabelSelector: labelSelector,
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return deploymentConfigs.List(options)
		},
	}).List
}

func imageStreamsLister(imageStreams client.ImageStreamInterface, resourcesChan chan<- openshift.Resource, labelSelector string) func() error {
	glog.V(1).Infof("Listing ImageStreams...")
	return (&openshift.ExportLister{
		ResourcesChan: resourcesChan,
		LabelSelector: labelSelector,
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return imageStreams.List(options)
		},
	}).List
}

func routesLister(routes client.RouteInterface, resourcesChan chan<- openshift.Resource, labelSelector string) func() error {
	glog.V(1).Infof("Listing Routes...")
	return (&openshift.ExportLister{
		ResourcesChan: resourcesChan,
		LabelSelector: labelSelector,
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return routes.List(options)
		},
	}).List
}

func templatesLister(templates client.TemplateInterface, resourcesChan chan<- openshift.Resource, labelSelector string) func() error {
	glog.V(1).Infof("Listing Templates...")
	return (&openshift.ExportLister{
		ResourcesChan: resourcesChan,
		LabelSelector: labelSelector,
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return templates.List(options)
		},
	}).List
}

func usersLister(users client.UserInterface, resourcesChan chan<- openshift.Resource, labelSelector string) func() error {
	glog.V(1).Infof("Listing Users...")
	return (&openshift.ExportLister{
		ResourcesChan: resourcesChan,
		LabelSelector: labelSelector,
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return users.List(options)
		},
	}).List
}

func groupsLister(groups client.GroupInterface, resourcesChan chan<- openshift.Resource, labelSelector string) func() error {
	glog.V(1).Infof("Listing Groups...")
	return (&openshift.ExportLister{
		ResourcesChan: resourcesChan,
		LabelSelector: labelSelector,
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return groups.List(options)
		},
	}).List
}

func policiesLister(policies client.PolicyInterface, resourcesChan chan<- openshift.Resource, labelSelector string) func() error {
	glog.V(1).Infof("Listing Policies...")
	return (&openshift.ExportLister{
		ResourcesChan: resourcesChan,
		LabelSelector: labelSelector,
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return policies.List(options)
		},
	}).List
}

func policyBindingsLister(policyBindings client.PolicyBindingInterface, resourcesChan chan<- openshift.Resource, labelSelector string) func() error {
	glog.V(1).Infof("Listing PolicyBindings...")
	return (&openshift.ExportLister{
		ResourcesChan: resourcesChan,
		LabelSelector: labelSelector,
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return policyBindings.List(options)
		},
	}).List
}

func clusterPoliciesLister(clusterPolicies client.ClusterPolicyInterface, resourcesChan chan<- openshift.Resource, labelSelector string) func() error {
	glog.V(1).Infof("Listing ClusterPolicies...")
	return (&openshift.ExportLister{
		ResourcesChan: resourcesChan,
		LabelSelector: labelSelector,
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return clusterPolicies.List(options)
		},
	}).List
}

func clusterPolicyBindingsLister(clusterPolicyBindings client.ClusterPolicyBindingInterface, resourcesChan chan<- openshift.Resource, labelSelector string) func() error {
	glog.V(1).Infof("Listing ClusterPolicyBindings...")
	return (&openshift.ExportLister{
		ResourcesChan: resourcesChan,
		LabelSelector: labelSelector,
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return clusterPolicyBindings.List(options)
		},
	}).List
}

func replicationControllersLister(replicationControllers unversioned.ReplicationControllerInterface, resourcesChan chan<- openshift.Resource, labelSelector string) func() error {
	glog.V(1).Infof("Listing ReplicationControllers...")
	return (&openshift.ExportLister{
		ResourcesChan: resourcesChan,
		LabelSelector: labelSelector,
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return replicationControllers.List(options)
		},
		Requirements: []func() (*labels.Requirement, error){
			func() (*labels.Requirement, error) {
				return labels.NewRequirement(deployapi.DeploymentConfigAnnotation, labels.DoesNotExistOperator, sets.NewString())
			},
		},
	}).List
}

func podsLister(pods unversioned.PodInterface, resourcesChan chan<- openshift.Resource, labelSelector string) func() error {
	glog.V(1).Infof("Listing Pods...")
	return (&openshift.ExportLister{
		ResourcesChan: resourcesChan,
		LabelSelector: labelSelector,
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return pods.List(options)
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
	}).List
}

func endpointsLister(endpoints unversioned.EndpointsInterface, resourcesChan chan<- openshift.Resource, labelSelector string) func() error {
	glog.V(1).Infof("Listing Endpoints...")
	return (&openshift.ExportLister{
		ResourcesChan: resourcesChan,
		LabelSelector: labelSelector,
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return endpoints.List(options)
		},
	}).List
}

func servicesLister(services unversioned.ServiceInterface, resourcesChan chan<- openshift.Resource, labelSelector string) func() error {
	glog.V(1).Infof("Listing Services...")
	return (&openshift.ExportLister{
		ResourcesChan: resourcesChan,
		LabelSelector: labelSelector,
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return services.List(options)
		},
	}).List
}

func serviceAccountsLister(serviceAccounts unversioned.ServiceAccountsInterface, resourcesChan chan<- openshift.Resource, labelSelector string) func() error {
	glog.V(1).Infof("Listing ServiceAccounts...")
	return (&openshift.ExportLister{
		ResourcesChan: resourcesChan,
		LabelSelector: labelSelector,
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return serviceAccounts.List(options)
		},
	}).List
}

func secretsLister(secrets unversioned.SecretsInterface, resourcesChan chan<- openshift.Resource, labelSelector string) func() error {
	glog.V(1).Infof("Listing Secrets...")
	return (&openshift.ExportLister{
		ResourcesChan: resourcesChan,
		LabelSelector: labelSelector,
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return secrets.List(options)
		},
	}).List
}

func limitRangesLister(limitRanges unversioned.LimitRangeInterface, resourcesChan chan<- openshift.Resource, labelSelector string) func() error {
	glog.V(1).Infof("Listing LimitRanges...")
	return (&openshift.ExportLister{
		ResourcesChan: resourcesChan,
		LabelSelector: labelSelector,
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return limitRanges.List(options)
		},
	}).List
}

func resourceQuotasLister(resourceQuotas unversioned.ResourceQuotaInterface, resourcesChan chan<- openshift.Resource, labelSelector string) func() error {
	glog.V(1).Infof("Listing ResourceQuotas...")
	return (&openshift.ExportLister{
		ResourcesChan: resourcesChan,
		LabelSelector: labelSelector,
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return resourceQuotas.List(options)
		},
	}).List
}

func persistentVolumeClaimsLister(persistentVolumeClaims unversioned.PersistentVolumeClaimInterface, resourcesChan chan<- openshift.Resource, labelSelector string) func() error {
	glog.V(1).Infof("Listing PersistentVolumeClaims...")
	return (&openshift.ExportLister{
		ResourcesChan: resourcesChan,
		LabelSelector: labelSelector,
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return persistentVolumeClaims.List(options)
		},
	}).List
}

func persistentVolumesLister(persistentVolumes unversioned.PersistentVolumeInterface, resourcesChan chan<- openshift.Resource, labelSelector string) func() error {
	glog.V(1).Infof("Listing PersistentVolumes...")
	return (&openshift.ExportLister{
		ResourcesChan: resourcesChan,
		LabelSelector: labelSelector,
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return persistentVolumes.List(options)
		},
	}).List
}

func securityContextConstraintsLister(securityContextConstraints unversioned.SecurityContextConstraintInterface, resourcesChan chan<- openshift.Resource, labelSelector string) func() error {
	glog.V(1).Infof("Listing SecurityContextConstraints...")
	return (&openshift.ExportLister{
		ResourcesChan: resourcesChan,
		LabelSelector: labelSelector,
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return securityContextConstraints.List(options)
		},
	}).List
}

func namespacesLister(namespaces unversioned.NamespaceInterface, resourcesChan chan<- openshift.Resource, labelSelector string) func() error {
	glog.V(1).Infof("Listing Namespaces...")
	return (&openshift.ExportLister{
		ResourcesChan: resourcesChan,
		LabelSelector: labelSelector,
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return namespaces.List(options)
		},
	}).List
}

func namespaceLister(namespaces unversioned.NamespaceInterface, namespace string, resourcesChan chan<- openshift.Resource, labelSelector string) func() error {
	glog.V(1).Infof("Listing Namespace %s...", namespace)
	return (&openshift.ExportLister{
		ResourcesChan: resourcesChan,
		LabelSelector: labelSelector,
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
	}).List
}
