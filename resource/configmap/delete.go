package configmap

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Resource) ApplyDeleteChange(ctx context.Context, obj, deleteChange interface{}) error {
	configMapsToDelete, err := toConfigMaps(deleteChange)
	if err != nil {
		return microerror.Mask(err)
	}

	for _, configMap := range configMapsToDelete {
		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("deleting ConfigMap %#q in namespace %#q", configMap.Name, configMap.Namespace))

		err := r.k8sClient.CoreV1().ConfigMaps(configMap.Namespace).Delete(configMap.Name, &metav1.DeleteOptions{})
		if apierrors.IsNotFound(err) {
			r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("already deleted ConfigMap %#q in namespace %#q", configMap.Name, configMap.Namespace))
		} else if err != nil {
			return microerror.Mask(err)
		} else {
			r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("deleted ConfigMap %#q in namespace %#q", configMap.Name, configMap.Namespace))
		}

	}

	return nil
}

func (r *Resource) newDeleteChange(ctx context.Context, obj, currentState, desiredState interface{}) ([]*corev1.ConfigMap, error) {
	currentConfigMaps, err := toConfigMaps(currentState)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	desiredConfigMaps, err := toConfigMaps(desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var configMapsToDelete []*corev1.ConfigMap
	{
		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("computing ConfigMaps to delete"))

		for _, c := range currentConfigMaps {
			if !containsConfigMap(desiredConfigMaps, c) {
				configMapsToDelete = append(configMapsToDelete, c)
			}
		}

		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("computed %d ConfigMaps to delete", len(configMapsToDelete)))
	}

	return configMapsToDelete, nil
}
