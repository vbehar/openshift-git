package export

import (
	buildapi "github.com/openshift/origin/pkg/build/api"
	deployapi "github.com/openshift/origin/pkg/deploy/api"

	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/util/sets"
)

// DefaultRequirementsFor returns a list of requirements for the given kind
// that should be default requirements applied when listing/watching
// for example, ignore pods managed by RC and DC
func DefaultRequirementsFor(gvk unversioned.GroupVersionKind) []func() (*labels.Requirement, error) {
	switch gvk.Kind {

	case "ReplicationController":
		return []func() (*labels.Requirement, error){
			func() (*labels.Requirement, error) {
				return labels.NewRequirement(deployapi.DeploymentConfigAnnotation, labels.DoesNotExistOperator, sets.NewString())
			},
		}

	case "Pod":
		return []func() (*labels.Requirement, error){
			func() (*labels.Requirement, error) {
				return labels.NewRequirement(buildapi.BuildLabel, labels.DoesNotExistOperator, sets.NewString())
			},
			func() (*labels.Requirement, error) {
				return labels.NewRequirement(deployapi.DeployerPodForDeploymentLabel, labels.DoesNotExistOperator, sets.NewString())
			},
			func() (*labels.Requirement, error) {
				return labels.NewRequirement(deployapi.DeploymentConfigLabel, labels.DoesNotExistOperator, sets.NewString())
			},
		}
	}

	return nil
}
