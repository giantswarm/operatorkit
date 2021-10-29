package customresourcedefinition

import (
	"context"

	"github.com/giantswarm/microerror"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (w Wrapper) CreateObject(ctx context.Context, namespace string, obj interface{}) (interface{}, error) {
	configMap, err := toCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	createConfigMap, err := w.extClient.ApiextensionsV1().CustomResourceDefinitions().Create(ctx, &configMap, metav1.CreateOptions{})
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return createConfigMap, nil
}

func (w Wrapper) DeleteObject(ctx context.Context, name, namespace string) error {
	err := w.extClient.ApiextensionsV1().CustomResourceDefinitions().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (w Wrapper) GetObject(ctx context.Context, name, namespace string) (interface{}, error) {
	configMap, err := w.extClient.ApiextensionsV1().CustomResourceDefinitions().Get(ctx, name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return nil, microerror.Mask(notFoundError)
	} else if err != nil {
		return nil, microerror.Mask(err)
	}

	return configMap, nil
}

func (w Wrapper) UpdateObject(ctx context.Context, namespace string, obj interface{}) (interface{}, error) {
	configMap, err := toCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	m, err := w.extClient.ApiextensionsV1().CustomResourceDefinitions().Get(ctx, configMap.GetName(), metav1.GetOptions{})
	if err != nil {
		return nil, microerror.Mask(err)
	}
	configMap.SetResourceVersion(m.GetResourceVersion())

	updateConfigMap, err := w.extClient.ApiextensionsV1().CustomResourceDefinitions().Update(ctx, &configMap, metav1.UpdateOptions{})
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return updateConfigMap, nil
}

func toCustomObject(v interface{}) (apiextensionsv1.CustomResourceDefinition, error) {
	if v == nil {
		return apiextensionsv1.CustomResourceDefinition{}, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", &apiextensionsv1.CustomResourceDefinition{}, v)
	}

	customObjectPointer, ok := v.(*apiextensionsv1.CustomResourceDefinition)
	if !ok {
		return apiextensionsv1.CustomResourceDefinition{}, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", &apiextensionsv1.CustomResourceDefinition{}, v)
	}
	customObject := *customObjectPointer

	return customObject, nil
}
