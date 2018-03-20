// +build k8srequired

package integration

import (
	"github.com/giantswarm/microerror"
	corev1 "k8s.io/api/core/v1"
)

func createPod(pod *corev1.Pod) error {
	_, err := k8sClient.CoreV1().Pods(namespace).Create(pod)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func deletePod(name string) error {
	err := k8sClient.CoreV1().Pods(namespace).Delete(name, nil)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
