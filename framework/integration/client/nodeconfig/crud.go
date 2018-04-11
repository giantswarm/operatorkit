// +build k8srequired

package nodeconfig

import (
	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/microerror"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c Client) CreateObject(namespace string, obj interface{}) (interface{}, error) {
	nodeConfig, err := toCustomObject(obj)
	if err != nil {
		return nil, err
	}
	createNodeConfig, err := c.g8sClient.CoreV1alpha1().NodeConfigs(namespace).Create(&nodeConfig)
	if err != nil {
		return nil, err
	}

	return createNodeConfig, nil
}

func (c Client) DeleteObject(name, namespace string) error {
	err := c.g8sClient.CoreV1alpha1().NodeConfigs(namespace).Delete(name, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c Client) GetObject(name, namespace string) (interface{}, error) {
	nodeConfig, err := c.g8sClient.CoreV1alpha1().NodeConfigs(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return nodeConfig, nil
}

func (c Client) UpdateObject(namespace string, obj interface{}) (interface{}, error) {
	return nil, nil
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
