package drainerconfig

import (
	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/microerror"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (w Wrapper) CreateObject(namespace string, obj interface{}) (interface{}, error) {
	nodeConfig, err := toCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	createDrainerConfig, err := w.g8sClient.CoreV1alpha1().DrainerConfigs(namespace).Create(&nodeConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return createDrainerConfig, nil
}

func (w Wrapper) DeleteObject(name, namespace string) error {
	err := w.g8sClient.CoreV1alpha1().DrainerConfigs(namespace).Delete(name, nil)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (w Wrapper) GetObject(name, namespace string) (interface{}, error) {
	nodeConfig, err := w.g8sClient.CoreV1alpha1().DrainerConfigs(namespace).Get(name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return nil, microerror.Mask(notFoundError)
	} else if err != nil {
		return nil, microerror.Mask(err)
	}

	return nodeConfig, nil
}

func (w Wrapper) UpdateObject(namespace string, obj interface{}) (interface{}, error) {
	nodeConfig, err := toCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	m, err := w.g8sClient.CoreV1alpha1().DrainerConfigs(namespace).Get(nodeConfig.GetName(), metav1.GetOptions{})
	if err != nil {
		return nil, microerror.Mask(err)
	}
	nodeConfig.SetResourceVersion(m.GetResourceVersion())

	updateDrainerConfig, err := w.g8sClient.CoreV1alpha1().DrainerConfigs(namespace).Update(&nodeConfig)
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
