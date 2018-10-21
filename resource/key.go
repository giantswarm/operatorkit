package resource

import (
	"strings"

	"github.com/giantswarm/microerror"
	"k8s.io/apimachinery/pkg/api/meta"
)

const (
	// labelObject is a label added to all objects created by resoucres of
	// this package. It is used for deletion purposes. The value of the
	// label is namespace and name of watched object in format
	// "${namespace}-${name}". Together with labelResource it allows to
	// find objects created by a resource.
	labelObject = "operatorkit.giantswarm.io/object"
	// labelResource is label added to all objects created by resources of
	// this package. It is used for deletion purposes. The value of the
	// label is the Name property of the setting resource. Together with
	// labelObject it allows to find objects created by a resource.
	labelResource = "operatorkit.giantswarm.io/resource"
)

type assertLabelsType interface {
	printNameType
	GetLabels() map[string]string
}

func assertLabels(obj assertLabelsType, watchedObj interface{}, resourceName string) error {
	desired, err := newLabels(watchedObj, resourceName)
	if err != nil {
		return microerror.Mask(err)
	}

	for k, v := range desired {
		current := obj.GetLabels()
		if current == nil {
			return microerror.Maskf(executionFailedError, " %#q object %#q label is not set", newPrintName(obj), k)
		}
		value, ok := current[k]
		if !ok {
			return microerror.Maskf(executionFailedError, " %#q object %#q label is not set", newPrintName(obj), k)
		}
		if value != v {
			return microerror.Maskf(executionFailedError, " %#q object %#q label has value %#q but want %#q", newPrintName(obj), k, value, v)
		}
	}

	return nil
}

func newLabels(watchedObj interface{}, resourceName string) (map[string]string, error) {
	var labelObjectValue string
	{
		accessor, err := meta.Accessor(watchedObj)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		n := accessor.GetName()
		ns := accessor.GetNamespace()
		if ns == "" {
			ns = "default"
		}

		labelObjectValue = ns + "-" + n
	}

	labels := map[string]string{
		labelObject:   labelObjectValue,
		labelResource: resourceName,
	}

	return labels, nil
}

func newLabelSelector(watchedObj interface{}, resourceName string) (string, error) {
	labels, err := newLabels(watchedObj, resourceName)
	if err != nil {
		return "", microerror.Mask(err)
	}

	var selectors []string
	for k, v := range labels {
		s := k + "=" + v
		selectors = append(selectors, s)
	}

	return strings.Join(selectors, ","), nil
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

type setLabelsType interface {
	printNameType
	GetLabels() map[string]string
	SetLabels(map[string]string)
}

func setLabels(obj setLabelsType, watchedObj interface{}, resourceName string) error {
	desired, err := newLabels(watchedObj, resourceName)
	if err != nil {
		return microerror.Mask(err)
	}

	current := obj.GetLabels()

	if current == nil {
		obj.SetLabels(desired)
		return nil
	}

	for k, v := range desired {
		_, ok := current[k]
		if ok {
			return microerror.Maskf(executionFailedError, " %#q object %#q label already exists", newPrintName(obj), k)
		}
		current[k] = v
	}

	obj.SetLabels(current)
	return nil
}
