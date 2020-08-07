package drainerconfig

import (
	"context"

	"github.com/giantswarm/apiextensions/v2/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/microerror"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (w Wrapper) CreateObject(ctx context.Context, namespace string, obj interface{}) (interface{}, error) {
	drainerConfig, err := toCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	createDrainerConfig, err := w.k8sClient.G8sClient().CoreV1alpha1().DrainerConfigs(namespace).Create(ctx, &drainerConfig, metav1.CreateOptions{})
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return createDrainerConfig, nil
}

func (w Wrapper) DeleteObject(ctx context.Context, name, namespace string) error {
	err := w.k8sClient.G8sClient().CoreV1alpha1().DrainerConfigs(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (w Wrapper) GetObject(ctx context.Context, name, namespace string) (interface{}, error) {
	drainerConfig, err := w.k8sClient.G8sClient().CoreV1alpha1().DrainerConfigs(namespace).Get(ctx, name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return nil, microerror.Mask(notFoundError)
	} else if err != nil {
		return nil, microerror.Mask(err)
	}

	return drainerConfig, nil
}

func (w Wrapper) UpdateObject(ctx context.Context, namespace string, obj interface{}) (interface{}, error) {
	drainerConfig, err := toCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	m, err := w.k8sClient.G8sClient().CoreV1alpha1().DrainerConfigs(namespace).Get(ctx, drainerConfig.GetName(), metav1.GetOptions{})
	if err != nil {
		return nil, microerror.Mask(err)
	}
	drainerConfig.SetResourceVersion(m.GetResourceVersion())

	updateDrainerConfig, err := w.k8sClient.G8sClient().CoreV1alpha1().DrainerConfigs(namespace).Update(ctx, &drainerConfig, metav1.UpdateOptions{})
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return updateDrainerConfig, nil
}

func toCustomObject(v interface{}) (v1alpha1.DrainerConfig, error) {
	if v == nil {
		return v1alpha1.DrainerConfig{}, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", &v1alpha1.DrainerConfig{}, v)
	}

	customObjectPointer, ok := v.(*v1alpha1.DrainerConfig)
	if !ok {
		return v1alpha1.DrainerConfig{}, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", &v1alpha1.DrainerConfig{}, v)
	}
	customObject := *customObjectPointer

	return customObject, nil
}
