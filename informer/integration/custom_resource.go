// +build k8srequired

package informer

import (
	"testing"

	"github.com/giantswarm/microerror"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
)

const (
	namespace = "informer-integration-test"
)

type CustomResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              CustomResourceSpec `json:"spec"`
}

func (t *CustomResource) DeepCopyObject() runtime.Object {
	return &CustomResource{
		TypeMeta:   t.TypeMeta,
		ObjectMeta: *t.ObjectMeta.DeepCopy(),
		Spec:       t.Spec,
	}
}

type CustomResourceSpec struct {
	ID string `json:"id" yaml:"id"`
}

type CustomResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []CustomResource `json:"items"`
}

func (t *CustomResourceList) DeepCopyObject() runtime.Object {
	itemsCopy := make([]CustomResource, len(t.Items))
	for i, item := range t.Items {
		itemsCopy[i] = item
	}

	return &CustomResourceList{
		TypeMeta: t.TypeMeta,
		ListMeta: *t.ListMeta.DeepCopy(),
		Items:    itemsCopy,
	}
}

func mustAssertCRWithID(e watch.Event, IDs ...string) {
	m, err := meta.Accessor(e.Object)
	if err != nil {
		panic(err)
	}

	name := m.GetName()
	for _, ID := range IDs {
		if name == ID {
			return
		}
	}

	panic("expected one of %#v got %#v", IDs, name)
}

func createCustomResource(t *testing.T, ID string) error {
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

func deleteCustomResource(t *testing.T, ID string) {
	err := k8sClient.CoreV1().ConfigMaps(namespace).Delete(ID, nil)
	if err != nil {
		t.Fatalf("expected %#v got %#v", nil, err)
	}
}
