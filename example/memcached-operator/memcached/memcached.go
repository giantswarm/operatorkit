package memcached

import (
	"github.com/giantswarm/apiextensions/pkg/clientset/versioned"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/client/k8srestconfig"
	"github.com/giantswarm/operatorkit/controller"
	"k8s.io/client-go/rest"
)

// Config represents the configuration used to create a new memcached controller.
type Config struct {
	Logger       micrologger.Logger
	K8sAddress   string
	K8sInCluster bool
	K8sCrtFile   string
	K8sKeyFile   string
	K8sCAFile    string
}

// Memcached is a type containing the OperatorKit controller.
type Memcached struct {
	*controller.Controller
}

// New creates a new memcached controller.
func New(config Config) (*Memcached, error) {
	var err error

	var restConfig *rest.Config
	{
		c := k8srestconfig.Config{
			Logger: config.Logger,

			Address:   config.K8sAddress,
			InCluster: config.K8sInCluster,
			TLS: k8srestconfig.TLSClientConfig{
				CAFile:  config.K8sCAFile,
				CrtFile: config.K8sCrtFile,
				KeyFile: config.K8sKeyFile,
			},
		}

		restConfig, err = k8srestconfig.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	_, err = versioned.NewForConfig(restConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return &Memcached{}, nil
}
