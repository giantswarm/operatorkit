package configmapresource

import (
	"context"

	"github.com/giantswarm/microerror"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Resource) ApplyUpdateChange(ctx context.Context, obj, updateChange interface{}) error {
	configMaps, err := toConfigMaps(updateChange)
	if err != nil {
		return microerror.Mask(err)
	}

	for _, configMap := range configMaps {
		r.logger.Debugf(ctx, "updating ConfigMap %#q in namespace %#q", configMap.Name, configMap.Namespace)

		_, err = r.k8sClient.CoreV1().ConfigMaps(configMap.Namespace).Update(ctx, configMap, metav1.UpdateOptions{})
		if err != nil {
			return microerror.Mask(err)
		}

		r.logger.Debugf(ctx, "updated ConfigMap %#q in namespace %#q", configMap.Name, configMap.Namespace)
	}

	return nil
}

func (r *Resource) newUpdateChange(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	currentConfigMaps, err := toConfigMaps(currentState)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	desiredConfigMaps, err := toConfigMaps(desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var configMapsToUpdate []*corev1.ConfigMap
	{
		r.logger.Debugf(ctx, "computing ConfigMaps to update")

		for _, c := range currentConfigMaps {
			for _, d := range desiredConfigMaps {
				m := newConfigMapToUpdate(c, d, r.allowedLabels)
				if m != nil {
					configMapsToUpdate = append(configMapsToUpdate, m)
				}
			}
		}

		r.logger.Debugf(ctx, "computed %d ConfigMaps to update", len(configMapsToUpdate))
	}

	return configMapsToUpdate, nil
}
