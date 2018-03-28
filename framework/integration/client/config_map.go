// +build k8srequired

package client

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateConfigMap(configMap *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	createConfigMap, err := k8sClient.CoreV1().ConfigMaps(configMap.Namespace).Create(configMap)
	if err != nil {
		return nil, err
	}

	return createConfigMap, nil
}

func GetConfigMap(name, namespace string) (*corev1.ConfigMap, error) {
	configMap, err := k8sClient.CoreV1().ConfigMaps(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return configMap, nil
}

func DeleteConfigMap(name, namespace string) error {
	err := k8sClient.CoreV1().ConfigMaps(namespace).Delete(name, nil)
	if err != nil {
		return err
	}

	return nil
}
