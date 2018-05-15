package resource

import (
	"context"
	"fmt"

	examplev1alpha1 "github.com/giantswarm/apiextensions/pkg/apis/example/v1alpha1"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/example/memcached-operator/controller/key"
	"github.com/giantswarm/operatorkit/example/memcached-operator/logger"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

const (
	servicesName = "services"
)

type ServicesConfig struct {
	K8sClient kubernetes.Interface
}

type Services struct {
	k8sClient kubernetes.Interface
}

func NewServices(config ServicesConfig) (*Services, error) {
	d := &Services{
		k8sClient: config.K8sClient,
	}

	return d, nil
}
func (s *Services) Name() string {
	return servicesName
}

func (s *Services) EnsureCreated(ctx context.Context, obj interface{}) error {
	memcachedConfig := obj.(*examplev1alpha1.MemcachedConfig).DeepCopy()

	managedServices, err := s.k8sClient.CoreV1().Services(memcachedConfig.Namespace).List(metav1.ListOptions{
		LabelSelector: key.LabelSelectorManagedBy,
	})
	if err != nil {
		return microerror.Mask(err)
	}

	currentReplicas := len(managedServices.Items)
	desiredReplicas := memcachedConfig.Spec.Replicas

	logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("currentReplicas = %d desiredReplicas = %d", currentReplicas, desiredReplicas))

	// Update existing services and scale up if necessary.
	for i := 0; i < desiredReplicas; i++ {
		err := s.ensureReplicaCreated(ctx, memcachedConfig, i)
		if err != nil {
			return microerror.Mask(err)
		}

	}

	// Scale down if necessary.
	for i := desiredReplicas; i < currentReplicas; i++ {
		err := s.ensureReplicaDeleted(ctx, memcachedConfig, i)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func (s *Services) EnsureDeleted(ctx context.Context, obj interface{}) error {
	memcachedConfig := obj.(*examplev1alpha1.MemcachedConfig).DeepCopy()

	managedServices, err := s.k8sClient.CoreV1().Services(memcachedConfig.Namespace).List(metav1.ListOptions{
		LabelSelector: key.LabelSelectorManagedBy,
	})
	if err != nil {
		return microerror.Mask(err)
	}

	currentReplicas := len(managedServices.Items)

	logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("currentReplicas = %d", currentReplicas))

	if currentReplicas == 0 {
		logger.LogCtx(ctx, "level", "debug", "message", "no created services")
		return nil
	}

	// Delete existing servicees.
	for i := 0; i < currentReplicas; i++ {
		err := s.ensureReplicaDeleted(ctx, memcachedConfig, i)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func (s *Services) ensureReplicaCreated(ctx context.Context, m *examplev1alpha1.MemcachedConfig, replica int) error {
	desired, err := newDesiredService(m, replica)
	if err != nil {
		return microerror.Mask(err)
	}

	current, err := s.k8sClient.CoreV1().Services(desired.Namespace).Get(desired.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		// Just make sure current is nil when not found.
		current = nil
	} else if err != nil {
		return microerror.Mask(err)
	}

	if current == nil {
		_, err = s.k8sClient.CoreV1().Services(desired.Namespace).Create(desired)
		if err != nil {
			return microerror.Mask(err)
		}

		logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("service %s/%s created", desired.Namespace, desired.Name))
	} else {
		// Service.spec.clusterIP is immutable and filled by the
		// Kubernetes when not provided.
		desired.Spec.ClusterIP = current.Spec.ClusterIP

		desired.ResourceVersion = current.ResourceVersion

		_, err = s.k8sClient.CoreV1().Services(desired.Namespace).Update(desired)
		if err != nil {
			return microerror.Mask(err)
		}

		logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("service %s/%s updated", desired.Namespace, desired.Name))
	}

	return nil
}

func (s *Services) ensureReplicaDeleted(ctx context.Context, m *examplev1alpha1.MemcachedConfig, replica int) error {
	name := key.ReplicaName(replica)
	namespace := key.Namespace(m)

	err := s.k8sClient.CoreV1().Services(namespace).Delete(name, &metav1.DeleteOptions{})
	if apierrors.IsNotFound(err) {
		logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("service %s/%s already deleted", namespace, name))
	} else if err != nil {
		return microerror.Mask(err)
	} else {
		logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("service %s/%s deleted", namespace, name))
	}

	return nil
}

func newDesiredService(memcachedConfig *examplev1alpha1.MemcachedConfig, replica int) (*corev1.Service, error) {
	name := key.ReplicaName(replica)
	namespace := key.Namespace(memcachedConfig)

	labels := key.CommonLabels(memcachedConfig, replica)

	s := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       key.MemcachedPort,
					TargetPort: intstr.FromInt(key.MemcachedPort),
				},
			},
		},
	}

	return s, nil
}
