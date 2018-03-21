// +build k8srequired

package integration

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createConfigMap(configMap *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	createConfigMap, err := k8sClient.CoreV1().ConfigMaps(namespace).Create(configMap)
	if err != nil {
		return nil, err
	}

	return createConfigMap, nil
}

func getConfigMap(name string) (*corev1.ConfigMap, error) {
	configMap, err := k8sClient.CoreV1().ConfigMaps(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return configMap, nil
}

func deleteConfigMap(name string) error {
	err := k8sClient.CoreV1().ConfigMaps(namespace).Delete(name, nil)
	if err != nil {
		return err
	}

	return nil
}
