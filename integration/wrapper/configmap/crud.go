package configmap

import (
	"context"

	"github.com/giantswarm/microerror"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (w Wrapper) CreateObject(ctx context.Context, namespace string, obj interface{}) (interface{}, error) {
	configMap, err := toCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	createConfigMap, err := w.k8sClient.CoreV1().ConfigMaps(namespace).Create(ctx, &configMap, metav1.CreateOptions{})
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return createConfigMap, nil
}

func (w Wrapper) DeleteObject(ctx context.Context, name, namespace string) error {
	err := w.k8sClient.CoreV1().ConfigMaps(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (w Wrapper) GetObject(ctx context.Context, name, namespace string) (interface{}, error) {
	configMap, err := w.k8sClient.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
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

	m, err := w.k8sClient.CoreV1().ConfigMaps(namespace).Get(ctx, configMap.GetName(), metav1.GetOptions{})
	if err != nil {
		return nil, microerror.Mask(err)
	}
	configMap.SetResourceVersion(m.GetResourceVersion())

	updateConfigMap, err := w.k8sClient.CoreV1().ConfigMaps(configMap.Namespace).Update(ctx, &configMap, metav1.UpdateOptions{})
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return updateConfigMap, nil
}

func toCustomObject(v interface{}) (corev1.ConfigMap, error) {
	if v == nil {
		return corev1.ConfigMap{}, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", &corev1.ConfigMap{}, v)
	}

	customObjectPointer, ok := v.(*corev1.ConfigMap)
	if !ok {
		return corev1.ConfigMap{}, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", &corev1.ConfigMap{}, v)
	}
	customObject := *customObjectPointer

	return customObject, nil
}
