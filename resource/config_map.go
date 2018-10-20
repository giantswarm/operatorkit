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
	// DesiredFn is function returning a desired ConfigMap objects for the
	// given custom resource object.
	DesiredFn func(context.Context, interface{}) ([]corev1.ConfigMap, error)
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger

	// DeleteOptions to use when deleting created objects. It is useful for
	// setting a grace period. This setting is optional.
	DeleteOptions metav1.DeleteOptions
	// Name is a name returned by Name method of the resource.
	Name string
}

type ConfigMap struct {
	desiredFn func(context.Context, interface{}) ([]corev1.ConfigMap, error)
	k8sClient kubernetes.Interface
	logger    micrologger.Logger

	deleteOptions metav1.DeleteOptions
	name          string
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

	c := &ConfigMap{
		desiredFn: config.DesiredFn,
		k8sClient: config.K8sClient,
		logger:    config.Logger,

		deleteOptions: config.DeleteOptions,
		name:          config.Name,
	}

	return c, nil
}

func (c *ConfigMap) Name() string {
	return c.name
}

func (c *ConfigMap) EnsureCreated(ctx context.Context, obj interface{}) error {
	var err error

	var desired []corev1.ConfigMap
	{
		c.logger.LogCtx(ctx, "level", "debug", "message", "finding desired ConfigMap objects")

		desired, err = c.desiredFn(ctx, obj)
		if err != nil {
			return microerror.Mask(err)
		}

		// Enrich desired objects with annotation to be able to select
		// and clean them during the deletion.
		for _, d := range desired {
			err = setAnnotation(&d, annotationResource, valueConfigMap)
			if err != nil {
				return microerror.Mask(err)
			}
		}

		c.logger.LogCtx(ctx, "level", "debug", "message", "found desired ConfigMap objects")
	}

	for _, d := range desired {
		err := c.ensureCreated(ctx, &d)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func (c *ConfigMap) ensureCreated(ctx context.Context, desired *corev1.ConfigMap) error {
	var err error

	printName := newPrintName(desired)

	var current *corev1.ConfigMap
	{
		c.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("finding current ConfigMap %#q", printName))

		opts := metav1.GetOptions{
			IncludeUninitialized: true,
		}

		current, err = c.k8sClient.CoreV1().ConfigMaps(desired.Namespace).Get(desired.Name, opts)
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

		// Make sure the object is managed by this resource by checking
		// annotations.
		err = assertAnnoatation(current, annotationResource, valueConfigMap)
		if err != nil {
			return microerror.Mask(err)
		}

		toUpdate := newConfigMapToUpdate(current, desired)

		if toUpdate != nil {
			_, err := c.k8sClient.CoreV1().ConfigMaps(toUpdate.Namespace).Update(toUpdate)
			if err != nil {
				return microerror.Mask(err)
			}

			c.logger.LogCtx(ctx, "level", "info", "message", fmt.Sprintf("updated ConfigMap %#q", printName))

		} else {
			c.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("ConfigMap %#q is up to date", printName))
		}

	} else {
		c.logger.LogCtx(ctx, "level", "info", "message", fmt.Sprintf("creating ConfigMap %#q", printName))

		_, err := c.k8sClient.CoreV1().ConfigMaps(desired.Namespace).Create(desired)
		if err != nil {
			return microerror.Mask(err)
		}

		c.logger.LogCtx(ctx, "level", "info", "message", fmt.Sprintf("created ConfigMap %#q", printName))
	}

	return nil
}

func (c *ConfigMap) EnsureDeleted(ctx context.Context, obj interface{}) error {
	var err error

	labelSelector := annotationResource + "=" + valueConfigMap

	var currentObjects []corev1.ConfigMap
	{
		c.logger.LogCtx(ctx, "level", "info", "message", fmt.Sprintf("finding ConfigMap objects for selector %#q", labelSelector))

		opts := metav1.ListOptions{
			LabelSelector: labelSelector,
			// We don't want orphaned uninitialized objects.
			IncludeUninitialized: true,
		}

		list, err := c.k8sClient.CoreV1().ConfigMaps("").List(opts)
		if err != nil {
			return microerror.Mask(err)
		}

		currentObjects = list.Items

		if len(currentObjects) == 0 {
			c.logger.LogCtx(ctx, "level", "info", "message", fmt.Sprintf("did not find ConfigMap objects for selector %#q", labelSelector))
		} else {
			c.logger.LogCtx(ctx, "level", "info", "message", fmt.Sprintf("found ConfigMap objects for selector %#q", labelSelector))
		}
	}

	for _, current := range currentObjects {
		printName := newPrintName(&current)

		c.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("deleting ConfigMap %#q", printName))

		err = c.k8sClient.CoreV1().ConfigMaps(current.Namespace).Delete(current.Name, &c.deleteOptions)
		if err != nil {
			return microerror.Mask(err)
		}

		c.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("deleting ConfigMap %#q", printName))
	}

	return nil
}

// newConfigMapToUpdate creates a new ConfigMap to update based on current and
// desired one. If there is no update nil is returned.
func newConfigMapToUpdate(current, desired *corev1.ConfigMap) *corev1.ConfigMap {
	merged := current.DeepCopy()

	if !reflect.DeepEqual(current.Annotations, desired.Annotations) {
		merged.Annotations = desired.Annotations
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
