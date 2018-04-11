// +build k8srequired

package configmap

import (
	"github.com/giantswarm/microerror"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c Client) CreateObject(namespace string, obj interface{}) (interface{}, error) {
	configMap, err := toCustomObject(obj)
	if err != nil {
		return nil, err
	}
	createConfigMap, err := c.k8sClient.CoreV1().ConfigMaps(namespace).Create(&configMap)
	if err != nil {
		return nil, err
	}

	return createConfigMap, nil
}

func (c Client) DeleteObject(name, namespace string) error {
	err := c.k8sClient.CoreV1().ConfigMaps(namespace).Delete(name, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c Client) GetObject(name, namespace string) (interface{}, error) {
	configMap, err := c.k8sClient.CoreV1().ConfigMaps(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return configMap, nil
}

func (c Client) UpdateObject(namespace string, obj interface{}) (interface{}, error) {
	return nil, nil
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
