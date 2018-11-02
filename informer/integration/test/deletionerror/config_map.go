// +build k8srequired

package deletionerror

import (
	"fmt"

	"github.com/giantswarm/microerror"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

func createConfigMap(ID string) error {
	cm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      ID,
			Namespace: namespace,
		},
		Data: map[string]string{},
	}

	_, err := k8sClient.CoreV1().ConfigMaps(namespace).Create(cm)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func deleteConfigMap(ID string) error {
	err := k8sClient.CoreV1().ConfigMaps(namespace).Delete(ID, nil)
	if err != nil {
		return microerror.Mask(err)
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
