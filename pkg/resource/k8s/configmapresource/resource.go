package configmapresource

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type Config struct {
	K8sClient   kubernetes.Interface
	Logger      micrologger.Logger
	StateGetter StateGetter

	AllowedLabels []string
	Name          string
}

type Resource struct {
	k8sClient   kubernetes.Interface
	logger      micrologger.Logger
	stateGetter StateGetter

	allowedLabels map[string]bool
	name          string
}

func New(config Config) (*Resource, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.StateGetter == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.StateGetter must not be empty", config)
	}

	if config.Name == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Name must not be empty", config)
	}

	r := &Resource{
		k8sClient:   config.K8sClient,
		logger:      config.Logger,
		stateGetter: config.StateGetter,

		name: config.Name,
	}

	if config.AllowedLabels != nil {
		allowedLabels := map[string]bool{}
		{
			for _, label := range config.AllowedLabels {
				allowedLabels[label] = true
			}
		}

		r.allowedLabels = allowedLabels
	}

	return r, nil
}

func (r *Resource) Name() string {
	return r.name
}

func containsConfigMap(configMaps []*corev1.ConfigMap, configMap *corev1.ConfigMap) bool {
	for _, a := range configMaps {
		if configMap.Name == a.Name && configMap.Namespace == a.Namespace {
			return true
		}
	}

	return false
}

func toConfigMaps(v interface{}) ([]*corev1.ConfigMap, error) {
	x, ok := v.([]*corev1.ConfigMap)
	if !ok {
		return nil, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", x, v)
	}

	return x, nil
}
