package integration

import (
	"github.com/giantswarm/microerror"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createConfigMap(configMap *corev1.ConfigMap) error {
	_, err := k8sClient.CoreV1().ConfigMaps(namespace).Create(configMap)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func getConfigMap(name string) (*corev1.ConfigMap, error) {
	configMap, err := k8sClient.CoreV1().ConfigMaps(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return configMap, nil
}

func deleteConfigMap(ID string) error {
	err := k8sClient.CoreV1().ConfigMaps(namespace).Delete(ID, nil)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
