package openshift

import "k8s.io/kubernetes/pkg/labels"

// extendSelector extends the given labelSelector with the given requirements
func extendSelector(selector labels.Selector, requirements ...func() (*labels.Requirement, error)) (labels.Selector, error) {
	if selector == nil {
		selector = labels.Everything()
	}

	for _, fn := range requirements {
		r, err := fn()
		if err != nil {
			return nil, err
		}
		selector = selector.Add(*r)
	}

	return selector, nil
}
