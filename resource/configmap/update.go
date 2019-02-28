package configmap

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	v1 "k8s.io/api/core/v1"

	"github.com/giantswarm/release-operator/service/controller/key"
)

func (r *Resource) ApplyUpdateChange(ctx context.Context, obj, updateChange interface{}) error {
	configMap, err := key.ToConfigMap(updateChange)
	if err != nil {
		return microerror.Mask(err)
	}

	if configMap != nil {
		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("updating ConfigMap %#q in namespace %#q", configMap.Name, configMap.Namespace))

		_, err = r.g8sClient.ApplicationV1alpha1().Apps(configMap.Namespace).Update(configMap)
		if err != nil {
			return microerror.Mask(err)
		}

		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("updated ConfigMap %#q in namespace %#q", configMap.Name, configMap.Namespace))
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

	var configMapsToUpdate []*v1.ConfigMap
	{
		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("computing ConfigMaps to update"))

		for _, c := range currentConfigMaps {
			for _, d := range desiredConfigMaps {
				m := newConfigMapToUpdate(c, d)
				if m != nil {
					configMapsToUpdate = append(configMapsToUpdate, m)
				}
			}
		}

		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("computed %d ConfigMaps to update", configMapsToUpdate))
	}

	return configMapsToUpdate, nil
}
