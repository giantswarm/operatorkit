package resource

import "github.com/giantswarm/microerror"

func newPrintName(namespace, name string) string {
	if namespace == "" {
		return name
	}
	return namespace + "/" + name
}

type validateObj interface {
	GetName() string
	GetNamespace() string
}

func validateDesiredObject(obj validateObj, namespace, name string) error {
	if obj.GetName() != name {
		return microerror.Maskf(invalidDesiredSateError, "expected name %#q, got %#q", name, obj.GetName())
	}
	if obj.GetNamespace() != namespace {
		return microerror.Maskf(invalidDesiredSateError, "expected namespace %#q, got %#q", namespace, obj.GetNamespace())
	}

	return nil
}
