package example

import (
	"context"

	"github.com/giantswarm/microerror"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "github.com/giantswarm/operatorkit/v6/api/v1"
)

func (w Wrapper) CreateObject(ctx context.Context, namespace string, obj interface{}) (interface{}, error) {
	configMap, err := toCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	err = w.ctrlClient.Create(ctx, &configMap)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return &configMap, nil
}

func (w Wrapper) DeleteObject(ctx context.Context, name, namespace string) error {
	err := w.ctrlClient.Delete(ctx, &v1.Example{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	})
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (w Wrapper) GetObject(ctx context.Context, name, namespace string) (interface{}, error) {
	var obj v1.Example
	err := w.ctrlClient.Get(ctx, client.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}, &obj)
	if errors.IsNotFound(err) {
		return nil, microerror.Mask(notFoundError)
	} else if err != nil {
		return nil, microerror.Mask(err)
	}

	return &obj, nil
}

func (w Wrapper) UpdateObject(ctx context.Context, namespace string, obj interface{}) (interface{}, error) {
	customObject, err := toCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var existing v1.Example
	err = w.ctrlClient.Get(ctx, client.ObjectKey{
		Namespace: namespace,
		Name:      customObject.Name,
	}, &existing)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	customObject.SetResourceVersion(existing.GetResourceVersion())

	err = w.ctrlClient.Update(ctx, &customObject)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return &customObject, nil
}

func toCustomObject(v interface{}) (v1.Example, error) {
	if v == nil {
		return v1.Example{}, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", &apiextensionsv1.CustomResourceDefinition{}, v)
	}

	customObjectPointer, ok := v.(*v1.Example)
	if !ok {
		return v1.Example{}, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", &apiextensionsv1.CustomResourceDefinition{}, v)
	}
	customObject := *customObjectPointer

	return customObject, nil
}
