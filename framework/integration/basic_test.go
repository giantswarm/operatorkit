// +build k8srequired

package integration

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_Finalizer_Integration_Basic(t *testing.T) {
	mustSetup()
	defer mustTeardown()
	operatorName := "test-operator"
	podName := "testpod"
	operatorkitFramework, err := newFramework(operatorName)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
			SelfLink:  "/some/path",
		},
	}
	err = createPod(pod)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
	operatorkitFramework.UpdateFunc(pod, pod)

	resultPod, err := getPod(podName)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
	expectedFinalizers := []string{
		"operatorkit.giantswarm.io/test-operator",
	}

	if !reflect.DeepEqual(resultPod.GetFinalizers(), expectedFinalizers) {
		t.Fatalf("finalizers == %v, want %v", resultPod.GetFinalizers(), expectedFinalizers)
	}

}
