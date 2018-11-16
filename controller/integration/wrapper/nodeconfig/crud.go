// +build k8srequired

package nodeconfig

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
	createNodeConfig, err := w.g8sClient.CoreV1alpha1().NodeConfigs(namespace).Create(&nodeConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return createNodeConfig, nil
}

func (w Wrapper) DeleteObject(name, namespace string) error {
	err := w.g8sClient.CoreV1alpha1().NodeConfigs(namespace).Delete(name, nil)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (w Wrapper) GetObject(name, namespace string) (interface{}, error) {
	nodeConfig, err := w.g8sClient.CoreV1alpha1().NodeConfigs(namespace).Get(name, metav1.GetOptions{})
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

	m, err := w.g8sClient.CoreV1alpha1().NodeConfigs(namespace).Get(nodeConfig.GetName(), metav1.GetOptions{})
	if err != nil {
		return nil, microerror.Mask(err)
	}
	nodeConfig.SetResourceVersion(m.GetResourceVersion())

	updateNodeConfig, err := w.g8sClient.CoreV1alpha1().NodeConfigs(namespace).Update(&nodeConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return updateNodeConfig, nil
}

func toCustomObject(v interface{}) (v1alpha1.NodeConfig, error) {
	if v == nil {
		return v1alpha1.NodeConfig{}, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", &v1alpha1.NodeConfig{}, v)
	}

	customObjectPointer, ok := v.(*v1alpha1.NodeConfig)
	if !ok {
		return v1alpha1.NodeConfig{}, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", &v1alpha1.NodeConfig{}, v)
	}
	customObject := *customObjectPointer

	return customObject, nil
}
