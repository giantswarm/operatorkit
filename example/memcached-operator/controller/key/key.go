package key

import (
	"fmt"

	examplev1alpha1 "github.com/giantswarm/apiextensions/pkg/apis/example/v1alpha1"
	"github.com/giantswarm/microerror"
	"k8s.io/apimachinery/pkg/api/resource"
)

const (
	ControllerName = "memcached-operator"

	MemcachedPort    = 11211
	MemcachedVersion = "1.4.39-alpine"

	LabelSelectorManagedBy = labelManagedBy + "=" + ControllerName

	labelManagedBy = "giantswarm.io/managed-by"
	labelApp       = "app"
)

func CommonLabels(m *examplev1alpha1.MemcachedConfig, replica int) map[string]string {
	return map[string]string{
		labelApp:       ReplicaName(replica),
		labelManagedBy: ControllerName,
	}
}

func MemoryMB(m *examplev1alpha1.MemcachedConfig) (int64, error) {
	q, err := resource.ParseQuantity(m.Spec.Memory)
	if err != nil {
		return 0, microerror.Mask(err)
	}

	return q.ScaledValue(resource.Mega), nil
}

func Namespace(m *examplev1alpha1.MemcachedConfig) string {
	return m.Namespace
}

func ReplicaName(replica int) string {
	return fmt.Sprintf("memcached%04d", replica)
}
