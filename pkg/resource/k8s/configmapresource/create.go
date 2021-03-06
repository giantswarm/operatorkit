package configmapresource

import (
	"context"

	"github.com/giantswarm/microerror"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ApplyCreateChange ensures the ConfigMap is created in the k8s api.
func (r *Resource) ApplyCreateChange(ctx context.Context, obj, createChange interface{}) error {
	configMaps, err := toConfigMaps(createChange)
	if err != nil {
		return microerror.Mask(err)
	}

	for _, configMap := range configMaps {
		r.logger.Debugf(ctx, "creating ConfigMap %#q in namespace %#q", configMap.Name, configMap.Namespace)

		_, err = r.k8sClient.CoreV1().ConfigMaps(configMap.Namespace).Create(ctx, configMap, metav1.CreateOptions{})
		if apierrors.IsAlreadyExists(err) {
			r.logger.Debugf(ctx, "already created ConfigMap %#q in namespace %#q", configMap.Name, configMap.Namespace)
		} else if err != nil {
			return microerror.Mask(err)
		} else {
			r.logger.Debugf(ctx, "created ConfigMap %#q in namespace %#q", configMap.Name, configMap.Namespace)
		}
	}

	return nil
}

func (r *Resource) newCreateChange(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	currentConfigMaps, err := toConfigMaps(currentState)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	desiredConfigMaps, err := toConfigMaps(desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var configMapsToCreate []*corev1.ConfigMap
	{
		r.logger.Debugf(ctx, "computing ConfigMaps to create")

		for _, d := range desiredConfigMaps {
			if !containsConfigMap(currentConfigMaps, d) {
				configMapsToCreate = append(configMapsToCreate, d)
			}
		}

		r.logger.Debugf(ctx, "computed %d ConfigMaps to create", len(configMapsToCreate))
	}

	return configMapsToCreate, nil
}
