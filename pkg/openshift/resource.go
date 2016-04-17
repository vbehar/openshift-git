package openshift

import (
	"fmt"
	"strings"

	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/runtime"
)

// Resource represents an OpenShift resource
type Resource struct {
	// ObjectReference is the reference of the resource (kind, namespace, name, ...)
	*kapi.ObjectReference

	// Object is the resource itself
	Object runtime.Object

	// Exists indicates if the resource exists or not
	// (it may have been deleted)
	Exists bool

	// Status is a string representation of the current status of the resource
	// (like "added", "modified", "sync", or "deleted" for example)
	Status string
}

// NewResource instantiates a new Resource with its reference
// set to the given kind and namespacedName.
// The namespacedName is in the format "namespace/name" if it has a namespace
// or just "name" otherwise.
func NewResource(kind, namespacedName string) *Resource {
	var namespace, name string
	elems := strings.Split(namespacedName, "/")
	switch len(elems) {
	case 2:
		namespace = elems[0]
		name = elems[1]
	case 1:
		name = elems[0]
	}

	return &Resource{
		ObjectReference: &kapi.ObjectReference{
			Kind:      kind,
			Namespace: namespace,
			Name:      name,
		},
	}
}

// IsNamespaced returns true if the resource has a namespace
func (r *Resource) IsNamespaced() bool {
	return len(r.Namespace) > 0
}

// NamespacedName returns the namespacedName of the resource
// with the format "namespace/name" if it has a namespace
// or just "name" if it has no namespace
func (r *Resource) NamespacedName() string {
	if r.IsNamespaced() {
		return fmt.Sprintf("%s/%s", r.Namespace, r.Name)
	}
	return r.Name
}

// String returns a string representation of the resource
// (it's kind, namespace and name)
func (r *Resource) String() string {
	return fmt.Sprintf("%s %s", r.Kind, r.NamespacedName())
}
