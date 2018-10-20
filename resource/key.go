package resource

import "github.com/giantswarm/microerror"

const (
	// annotationResource is annotation added to all objects created by
	// resources of this package. It is used for deletion purposes. E.g.
	// all ConfigMap objects created by this package will have
	// "operatorkit.giantswarm.io/resource=ConfigMap" annotation set.
	annotationResource = "operatorkit.giantswarm.io/resource"

	valueConfigMap = "ConfigMap"
)

type assertAnnotationType interface {
	printNameType
	GetAnnotations() map[string]string
}

func assertAnnoatation(obj assertAnnotationType, key, value string) error {
	a := obj.GetAnnotations()
	if a == nil {
		return microerror.Maskf(executionFailedError, " %#q object %#q annotation is not set", newPrintName(obj), key)
	}
	v, ok := a[key]
	if !ok {
		return microerror.Maskf(executionFailedError, " %#q object %#q annotation is not set", newPrintName(obj), key)
	}
	if v != value {
		return microerror.Maskf(executionFailedError, " %#q object %#q annotation has value %#q but want %#q", newPrintName(obj), key, v, value)
	}

	return nil
}

type printNameType interface {
	GetName() string
	GetNamespace() string
}

func newPrintName(obj printNameType) string {
	if obj.GetNamespace() == "" {
		return obj.GetName()
	}
	return obj.GetNamespace() + "/" + obj.GetName()
}

type setAnnotationType interface {
	printNameType
	GetAnnotations() map[string]string
	SetAnnotations(map[string]string)
}

func setAnnotation(obj setAnnotationType, key, value string) error {
	a := obj.GetAnnotations()
	if a == nil {
		a = make(map[string]string, 1)
	}

	_, ok := a[key]
	if ok {
		return microerror.Maskf(executionFailedError, " %#q object %#q annotation already exists", newPrintName(obj), key)
	}

	a[key] = value
	obj.SetAnnotations(a)

	return nil
}
