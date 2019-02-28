package app

import (
	"github.com/giantswarm/apiextensions/pkg/apis/application/v1alpha1"
	"github.com/giantswarm/apiextensions/pkg/clientset/versioned"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"k8s.io/client-go/kubernetes"
)

type Config struct {
	G8sClient   versioned.Interface
	K8sClient   kubernetes.Interface
	Logger      micrologger.Logger
	StateGetter StateGetter

	Name string
}

type Resource struct {
	g8sClient   versioned.Interface
	k8sClient   kubernetes.Interface
	logger      micrologger.Logger
	stateGetter StateGetter

	name string
}

func New(config Config) (*Resource, error) {
	if config.G8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.G8sClient must not be empty", config)
	}
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
		g8sClient:   config.G8sClient,
		k8sClient:   config.K8sClient,
		logger:      config.Logger,
		stateGetter: config.StateGetter,

		name: config.Name,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return r.name
}

func containsAppCR(cr *v1alpha1.App, crs []*v1alpha1.App) bool {
	for _, a := range crs {
		if cr.Name == a.Name && cr.Namespace == a.Namespace {
			return true
		}
	}

	return false
}

func toAppCRs(v interface{}) ([]*v1alpha1.App, error) {
	x, ok := v.([]*v1alpha1.App)
	if !ok {
		return nil, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", x, v)
	}

	return x, nil
}
