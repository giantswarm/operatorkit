package k8sclient

import (
	"net/url"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	"github.com/giantswarm/microerror"
	"k8s.io/client-go/pkg/api"
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
	// Settings

	Address   string
	InCluster bool
	TLS       TLSClientConfig

	// Group is required only for NewRestClient.
	Group string
	// Version is required only for NewRestClient.
	Version string
}

// DefaultConfig provides a default configuration to create a new Kubernetes
// Clientset by best effort.
func DefaultConfig() Config {
	return Config{
		// Settings.

		Address:   "",
		InCluster: true,
		TLS:       TLSClientConfig{},
		Group:     "",
		Version:   "",
	}
}

func (c Config) Validate() error {
	// Settings.
	if c.Address == "" && !c.InCluster {
		return microerror.Maskf(invalidConfigError, "Address must not be empty when not creating in-cluster client")
	}

	if c.Address != "" {
		_, err := url.Parse(c.Address)
		if err != nil {
			return microerror.Maskf(invalidConfigError,
				"Address=%s must be a valid URL: %s", c.Address, err)
		}
	}

	return nil
}

func (c Config) ToK8sRestConfig() (*rest.Config, error) {
	var (
		config *rest.Config
		err    error
	)

	if c.InCluster {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, microerror.Mask(err)
		}
	} else {
		config = &rest.Config{
			Host: c.Address,
			TLSClientConfig: rest.TLSClientConfig{
				CertFile: c.TLS.CrtFile,
				KeyFile:  c.TLS.KeyFile,
				CAFile:   c.TLS.CAFile,
			},
		}
	}

	if c.Group != "" && c.Version != "" {
		config.GroupVersion = &schema.GroupVersion{
			Group:   c.Group,
			Version: c.Version,
		}
		config.NegotiatedSerializer = serializer.DirectCodecFactory{
			CodecFactory: api.Codecs,
		}
	}

	config.Burst = MaxBurst
	config.QPS = MaxQPS

	return config, nil
}
