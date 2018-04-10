// +build k8srequired

package client

import (
	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateNodeConfig(namespace string, nodeConfig *v1alpha1.NodeConfig) (*v1alpha1.NodeConfig, error) {
	createNodeConfig, err := g8sClient.CoreV1alpha1().NodeConfigs(namespace).Create(nodeConfig)
	if err != nil {
		return nil, err
	}

	return createNodeConfig, nil
}

func GetNodeConfig(name, namespace string) (*v1alpha1.NodeConfig, error) {
	nodeConfig, err := g8sClient.CoreV1alpha1().NodeConfigs(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return nodeConfig, nil
}

func DeleteNodeConfig(name, namespace string) error {
	err := g8sClient.CoreV1alpha1().NodeConfigs(namespace).Delete(name, nil)
	if err != nil {
		return err
	}

	return nil
}
