package k8sextclient

import (
	"net/url"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/rest"
)

const (
	// Maximum QPS to the master from this client.
	MaxQPS = 100
	// Maximum burst for throttle.
	MaxBurst = 100
)

// TLSClientConfig contains settings to enable transport layer security.
type TLSClientConfig struct {
	CAFile  string
	CrtFile string
	KeyFile string
}

// Config contains the common attributes to create a Kubernetes Clientset.
type Config struct {
	// Dependencies.
	Logger micrologger.Logger

	// Settings
	Address   string
	InCluster bool
	TLS       TLSClientConfig
}

// DefaultConfig provides a default configuration to create a new Kubernetes
// Clientset by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger: nil,

		// Settings.
		Address:   "",
		InCluster: true,
		TLS:       TLSClientConfig{},
	}
}

// New returns a Kubernetes extensions API Clientset with the provided
// configuration.
func New(config Config) (apiextensionsclient.Interface, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}

	// Settings.
	if config.Address == "" && !config.InCluster {
		return nil, microerror.Maskf(invalidConfigError, "config.Address must not be empty when not creating in-cluster client")
	}

	if config.Address != "" {
		_, err := url.Parse(config.Address)
		if err != nil {
			return nil, microerror.Maskf(invalidConfigError,
				"config.Address=%s must be a valid URL: %s", config.Address, err)
		}
	}

	var err error

	var restConfig *rest.Config
	if config.InCluster {
		config.Logger.Log("debug", "creating in-cluster config")

		restConfig, err = rest.InClusterConfig()
		if err != nil {
			return nil, microerror.Mask(err)
		}
	} else {
		config.Logger.Log("debug", "creating out-cluster config")

		restConfig = &rest.Config{
			Host: config.Address,
			TLSClientConfig: rest.TLSClientConfig{
				CertFile: config.TLS.CrtFile,
				KeyFile:  config.TLS.KeyFile,
				CAFile:   config.TLS.CAFile,
			},
		}
	}

	restConfig.Burst = MaxBurst
	restConfig.QPS = MaxQPS

	newClient, err := apiextensionsclient.NewForConfig(restConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return newClient, nil
}
