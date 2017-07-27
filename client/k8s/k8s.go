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
		Address:   "http://127.0.0.1:8080",
		InCluster: false,
		TLS:       TLSClientConfig{},
	}
}

// NewClient returns a Kubernetes Clientset with the provided configuration.
func NewClient(config Config) (kubernetes.Interface, error) {
	rawClientConfig, err := getRawClientConfig(config)
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(rawClientConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func newRawClientConfig(config Config) *rest.Config {
	tlsClientConfig := rest.TLSClientConfig{
		CertFile: config.TLS.CrtFile,
		KeyFile:  config.TLS.KeyFile,
		CAFile:   config.TLS.CAFile,
	}
	rawClientConfig := &rest.Config{
		Host:            config.Address,
		QPS:             MaxQPS,
		Burst:           MaxBurst,
		TLSClientConfig: tlsClientConfig,
	}

	return rawClientConfig
}

func getRawClientConfig(config Config) (*rest.Config, error) {
	var rawClientConfig *rest.Config
	var err error

	if config.InCluster {
		config.Logger.Log("debug", "creating in-cluster config")
		rawClientConfig, err = rest.InClusterConfig()
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	} else {
		if config.Address == "" {
			return nil, microerror.MaskAnyf(invalidConfigError, "kubernetes address must not be empty")
		}

		config.Logger.Log("debug", "creating out-cluster config")

		// Kubernetes listen URL.
		_, err := url.Parse(config.Address)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}

		rawClientConfig = newRawClientConfig(config)
	}

	return rawClientConfig, nil
}
