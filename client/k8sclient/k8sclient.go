package k8s

import (
	"net/url"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"k8s.io/client-go/kubernetes"
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
	var err error

	var newLogger micrologger.Logger
	{
		config := micrologger.DefaultConfig()
		newLogger, err = micrologger.New(config)
		if err != nil {
			panic(err)
		}
	}

	return Config{
		// Dependencies.
		Logger: newLogger,

		// Settings.
		Address:   "",
		InCluster: true,
		TLS:       TLSClientConfig{},
	}
}

// New returns a Kubernetes Clientset with the provided configuration.
func New(config Config) (kubernetes.Interface, error) {
	restConfig, err := toClientGoRESTConfig(config)
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func toClientGoRESTConfig(config Config) (*rest.Config, error) {
	if config.InCluster {
		config.Logger.Log("debug", "creating in-cluster config")
		c, err := rest.InClusterConfig()
		if err != nil {
			return nil, microerror.Mask(err)
		}
		return c, nil
	}

	if config.Address == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Address must not be empty")
	}
	_, err := url.Parse(config.Address)
	if err != nil {
		return nil, microerror.Maskf(invalidConfigError,
			"config.Address=%s must be a valid URL: %s", config.Address, err)
	}

	config.Logger.Log("debug", "creating out-cluster config")

	tlsClientConfig := rest.TLSClientConfig{
		CertFile: config.TLS.CrtFile,
		KeyFile:  config.TLS.KeyFile,
		CAFile:   config.TLS.CAFile,
	}

	c := &rest.Config{
		Host:            config.Address,
		QPS:             MaxQPS,
		Burst:           MaxBurst,
		TLSClientConfig: tlsClientConfig,
	}

	return c, nil
}
