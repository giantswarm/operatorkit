// +build k8srequired

package integration

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

func createConfigMap(namespace string, configMap *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	createConfigMap, err := k8sClient.CoreV1().ConfigMaps(namespace).Create(configMap)
	if err != nil {
		return nil, err
	}

	return createConfigMap, nil
}

func getConfigMap(namespace, name string) (*corev1.ConfigMap, error) {
	configMap, err := k8sClient.CoreV1().ConfigMaps(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return configMap, nil
}

func deleteConfigMap(namespace, name string) error {
	err := k8sClient.CoreV1().ConfigMaps(namespace).Delete(name, nil)
	if err != nil {
		return err
	}

	return nil
}

func mustAssertWithIDs(e watch.Event, IDs ...string) {
	configMap, ok := e.Object.(*corev1.ConfigMap)
	if !ok {
		panic(fmt.Sprintf("expected config map, got %#v", e.Object))
	}

	name := configMap.ObjectMeta.GetName()
	for _, ID := range IDs {
		if name == ID {
			return
		}
	}

	panic(fmt.Sprintf("expected one of %#v got %#v", IDs, name))
}
