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
func (d *Services) Name() string {
	return servicesName
}

func (d *Services) EnsureCreated(ctx context.Context, obj interface{}) error {
	memcachedConfig := obj.(*examplev1alpha1.MemcachedConfig).DeepCopy()

	managedServices, err := d.k8sClient.CoreV1().Services(memcachedConfig.Namespace).List(metav1.ListOptions{
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
		err := d.ensureReplicaCreated(ctx, memcachedConfig, i)
		if err != nil {
			return microerror.Mask(err)
		}

	}

	// Scale down if necessary.
	for i := desiredReplicas; i < currentReplicas; i++ {
		err := d.ensureReplicaDeleted(ctx, memcachedConfig, i)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func (d *Services) EnsureDeleted(ctx context.Context, obj interface{}) error {
	memcachedConfig := obj.(*examplev1alpha1.MemcachedConfig).DeepCopy()

	managedServices, err := d.k8sClient.CoreV1().Services(memcachedConfig.Namespace).List(metav1.ListOptions{
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
		err := d.ensureReplicaDeleted(ctx, memcachedConfig, i)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func (d *Services) ensureReplicaCreated(ctx context.Context, m *examplev1alpha1.MemcachedConfig, replica int) error {
	desired, err := newDesiredService(m, replica)
	if err != nil {
		return microerror.Mask(err)
	}

	_, err = d.k8sClient.CoreV1().Services(desired.Namespace).Update(desired)
	if apierrors.IsNotFound(err) {
		logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("service %s/%s does not exist, it will be created", desired.Namespace, desired.Name))
	} else if err != nil {
		return microerror.Mask(err)
	} else {
		logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("service %s/%s updated", desired.Namespace, desired.Name))
		return nil
	}

	_, err = d.k8sClient.CoreV1().Services(desired.Namespace).Create(desired)
	if apierrors.IsAlreadyExists(err) {
		logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("service %s/%s already exists", desired.Namespace, desired.Name))
	} else if err != nil {
		return microerror.Mask(err)
	} else {
		logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("service %s/%s created", desired.Namespace, desired.Name))
		return nil
	}

	return nil
}

func (d *Services) ensureReplicaDeleted(ctx context.Context, m *examplev1alpha1.MemcachedConfig, replica int) error {
	name := key.ReplicaName(replica)
	namespace := key.Namespace(m)

	err := d.k8sClient.CoreV1().Services(namespace).Delete(name, &metav1.DeleteOptions{})
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
