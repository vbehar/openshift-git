package openshift

import (
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/meta"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/util/sets"
)

var (
	// AllKinds is the list of all kinds that are likely candidates for exporting everything
	AllKinds = []string{
		// root-scoped kinds
		"ns", "pv", "scc", "clusterPolicies", "clusterPolicyBindings", "users", "groups",

		// namespace-scoped kinds
		"bc", "dc", "rc", "pods", "is", "svc", "routes", "templates", "secrets",
		"limits", "quota", "pvc", "policies", "policyBindings", "sa",
	}
)

// KindsFor parse the given list of kinds of resources (as string),
// and return a list of valid kinds (or an error).
// It supports the standard aliases, and our custom "everything" alias (see AllKinds).
func KindsFor(mapper meta.RESTMapper, kindsOrResources []string) ([]unversioned.GroupVersionKind, error) {
	resources := sets.NewString()
	for _, kindOrResource := range kindsOrResources {
		if kindOrResource == "everything" {
			for _, r := range AllKinds {
				resources.Insert(aliasesForResource(mapper, r)...)
			}
		} else {
			resources.Insert(aliasesForResource(mapper, kindOrResource)...)
		}
	}

	kindNames := sets.NewString()
	kinds := []unversioned.GroupVersionKind{}
	for _, resource := range resources.List() {
		gvr := kapi.Unversioned.WithResource(resource)
		gvk, err := mapper.KindFor(gvr)
		if err != nil {
			return []unversioned.GroupVersionKind{}, err
		}
		if !kindNames.Has(gvk.Kind) {
			kindNames.Insert(gvk.Kind)
			kinds = append(kinds, gvk)
		}
	}

	return kinds, nil
}

func aliasesForResource(mapper meta.RESTMapper, resource string) []string {
	if aliases, ok := mapper.AliasesForResource(resource); ok {
		return aliases
	}
	return []string{resource}
}
