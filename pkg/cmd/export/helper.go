package export

import (
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/meta"
	"k8s.io/kubernetes/pkg/kubectl"
	"k8s.io/kubernetes/pkg/runtime"
)

func upgradePrinterForObject(printer kubectl.ResourcePrinter, obj runtime.Object, mapper meta.RESTMapper) (kubectl.ResourcePrinter, error) {
	gvk, err := kapi.Scheme.ObjectKind(obj)
	if err != nil {
		return nil, err
	}

	mapping, err := mapper.RESTMapping(gvk.GroupKind())
	if err != nil {
		return nil, err
	}

	printer = kubectl.NewVersionedPrinter(printer, mapping.ObjectConvertor, mapping.GroupVersionKind.GroupVersion())
	return printer, nil
}
