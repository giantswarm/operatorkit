package resource

import (
	"context"
	"fmt"
	"reflect"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ConfigMapConfig struct {
	//
	// TODO describe & handle DesiredFn returning nil when it can't perpare desired state because there are some dependencies still missing.
	//
	// DesiredFn is function returning a desired ConfigMap object for the
	// given custom resource object. Returned ConfigMap must have name and
	// namespace equal to the values of Name and Namespace fields.
	DesiredFn func(context.Context, interface{}) (*corev1.ConfigMap, error)
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger

	// Name of the reconciled ConfigMap. It must be the same name as in
	// ConfigMap returned by DesiredFn.
	Name string
	// Namespace of the reconciled ConfigMap. It must be the same namespace
	// is in ConfigMap returned by DesiredFn.
	Namespace string
	// ResourceName is a name returned by Name method of the resource.
	ResourceName string
}

type ConfigMap struct {
	desiredFn func(context.Context, interface{}) (*corev1.ConfigMap, error)
	k8sClient kubernetes.Interface
	logger    micrologger.Logger

	name         string
	namespace    string
	resourceName string
}

func NewConfigMap(config ConfigMapConfig) (*ConfigMap, error) {
	if config.DesiredFn == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.DesiredFn must not be empty", config)
	}
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	if config.Name == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Name must not be empty", config)
	}
	if config.Namespace == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Namespace must not be empty", config)
	}
	if config.ResourceName == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.ResourceName must not be empty", config)
	}

	c := &ConfigMap{
		desiredFn: config.DesiredFn,
		k8sClient: config.K8sClient,
		logger:    config.Logger,

		name:         config.Name,
		namespace:    config.Namespace,
		resourceName: config.ResourceName,
	}

	return c, nil
}

func (c *ConfigMap) Name() string {
	return c.resourceName
}

func (c *ConfigMap) EnsureCreated(ctx context.Context, obj interface{}) error {
	var err error

	printName := newPrintName(c.namespace, c.name)

	var desired *corev1.ConfigMap
	{
		c.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("finding desired ConfigMap %#q", printName))

		desired, err = c.desiredFn(ctx, obj)
		if err != nil {
			return microerror.Mask(err)
		}

		validateDesiredObject(desired, c.namespace, c.name)
		if err != nil {
			return microerror.Mask(err)
		}

		c.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("found desired ConfigMap %#q", printName))
	}

	var current *corev1.ConfigMap
	{
		c.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("finding current ConfigMap %#q", printName))

		current, err = c.k8sClient.CoreV1().ConfigMaps(c.namespace).Get(c.name, metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			c.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("did not find current ConfigMap %#q", printName))

		} else if err != nil {
			return microerror.Mask(err)

		} else {
			c.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("found current ConfigMap %#q", printName))
		}
	}

	if current != nil {
		c.logger.LogCtx(ctx, "level", "info", "message", fmt.Sprintf("updating ConfigMap %#q", printName))

		toUpdate := newConfigMapToUpdate(current, desired)

		if toUpdate != nil {
			_, err := c.k8sClient.CoreV1().ConfigMaps(c.namespace).Update(toUpdate)
			if err != nil {
				return microerror.Mask(err)
			}

			c.logger.LogCtx(ctx, "level", "info", "message", fmt.Sprintf("updated ConfigMap %#q", printName))

		} else {
			c.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("ConfigMap %#q is up to date", printName))
		}

	} else {
		c.logger.LogCtx(ctx, "level", "info", "message", fmt.Sprintf("creating ConfigMap %#q", printName))

		_, err := c.k8sClient.CoreV1().ConfigMaps(c.namespace).Create(desired)
		if err != nil {
			return microerror.Mask(err)
		}

		c.logger.LogCtx(ctx, "level", "info", "message", fmt.Sprintf("created ConfigMap %#q", printName))
	}

	return nil
}

func (c *ConfigMap) EnsureDeleted(ctx context.Context, obj interface{}) error {
	var err error

	printName := newPrintName(c.namespace, c.name)

	{
		c.logger.LogCtx(ctx, "level", "info", "message", fmt.Sprintf("deleting ConfigMap %#q", printName))

		err = c.k8sClient.CoreV1().ConfigMaps(c.namespace).Delete(c.name, &metav1.DeleteOptions{})
		if apierrors.IsNotFound(err) {
			c.logger.LogCtx(ctx, "level", "info", "message", fmt.Sprintf("ConfigMap %#q already deleted", printName))

		} else if err != nil {
			return microerror.Mask(err)

		} else {
			c.logger.LogCtx(ctx, "level", "info", "message", fmt.Sprintf("deleted ConfigMap %#q", printName))
		}
	}

	return nil
}

// newConfigMapToUpdate creates a new ConfigMap to update based on current and
// desired one. If there is no update nil is returned.
func newConfigMapToUpdate(current, desired *corev1.ConfigMap) *corev1.ConfigMap {
	merged := current.DeepCopy()

	if !reflect.DeepEqual(current.Annotations, desired.Annotations) {
		merged.Labels = desired.Annotations
	}
	if !reflect.DeepEqual(current.Labels, desired.Labels) {
		merged.Labels = desired.Labels
	}

	if !reflect.DeepEqual(current.BinaryData, desired.BinaryData) {
		merged.BinaryData = desired.BinaryData
	}
	if !reflect.DeepEqual(current.Data, desired.Data) {
		merged.Data = desired.Data
	}

	if reflect.DeepEqual(current, merged) {
		return nil
	}

	return merged
}
