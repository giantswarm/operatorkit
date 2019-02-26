// Package k8srestconfig provides interface to create client-go rest config
// which can be used to construct various clients.
//
// Example usage:
//
//	import (
//		"k8s.io/client-go/kubernetes"
//		"k8s.io/client-go/rest"
//		apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
//
//		"github.com/giantswarm/operatorkit/client/k8srestconfig"
//		"github.com/giantswarm/apiextensions/pkg/clientset/versioned"
//		"github.com/giantswarm/microerror"
//	)
//
//	func f(config Config) error {
//		var err error
//
//		var restConfig *rest.Config
//		{
//			c := k8srestconfig.Config{
//				Logger: config.Logger,
//
//				Address:   config.Viper.GetString(config.Flag.Service.Kubernetes.Address),
//				InCluster: config.Viper.GetBool(config.Flag.Service.Kubernetes.InCluster),
//				TLS: TLSClientConfig{
//					CAFile:  config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.CAFile),
//					CrtFile: config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.CrtFile),
//					KeyFile: config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.KeyFile),
//				},
//			}
//
//			restConfig, err = k8srestconfig.New(c)
//			if err != nil {
//				return microerror.Mask(err)
//			}
//		}
//
//		k8sClient, err := kubernetes.NewForConfig(restConfig)
//		if err != nil {
//			return micorerror.Mask(err)
//		}
//
//		k8sExtClient, err := apiextensionsclient.NewForConfig(restConfig)
//		if err != nil {
//			return micorerror.Mask(err)
//		}
//
//		g8sClient, err := versioned.NewForConfig(restConfig)
//		if err != nil {
//			return microerror.Mask(err)
//		}
//	}
//
package k8srestconfig

import (
	"net/url"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"k8s.io/client-go/rest"
)

const (
	// Maximum QPS to the master from this client.
	MaxQPS = 100
	// Maximum burst for throttle.
	MaxBurst       = 100
	DefaultTimeout = 30 * time.Second
)

// TLSClientConfig contains settings to enable transport layer security.
type TLSClientConfig struct {
	// CAFile is the CA certificate for the cluster.
	CAFile string
	// CrtFile is the TLS client certificate.
	CrtFile string
	// KeyFile is the key for the TLS client certificate.
	KeyFile string
	// CAData holds PEM-encoded bytes. CAData takes precedence over CAFile.
	CAData []byte
	// CrtData holds PEM-encoded bytes. CrtData takes precedence over CrtFile.
	CrtData []byte
	// KeyData holds PEM-encoded bytes. KeyData takes precedence over KeyFile.
	KeyData []byte
}

// Config contains the common attributes to create a Kubernetes Clientset.
type Config struct {
	// Dependencies.
	Logger micrologger.Logger

	// Settings
	Address   string
	InCluster bool
	Timeout   time.Duration
	TLS       TLSClientConfig
}

// New returns a Kubernetes REST configuration for clients.
func New(config Config) (*rest.Config, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	// Settings.
	if config.Address == "" && !config.InCluster {
		return nil, microerror.Maskf(invalidConfigError, "%T.Address must not be empty when not creating in-cluster client", config)
	}

	if config.Address != "" {
		_, err := url.Parse(config.Address)
		if err != nil {
			return nil, microerror.Maskf(invalidConfigError, "%T.Address=%s must be a valid URL: %s", config, config.Address, err)
		}
	}
	if config.Timeout.Seconds() == 0 {
		config.Timeout = DefaultTimeout
	}

	var err error

	var restConfig *rest.Config
	if config.InCluster {
		config.Logger.Log("level", "debug", "message", "creating in-cluster config")

		restConfig, err = rest.InClusterConfig()
		if err != nil {
			return nil, microerror.Mask(err)
		}
	} else {
		config.Logger.Log("level", "debug", "message", "creating out-cluster config")

		restConfig = &rest.Config{
			Host:    config.Address,
			Timeout: config.Timeout,
			TLSClientConfig: rest.TLSClientConfig{
				CertFile: config.TLS.CrtFile,
				KeyFile:  config.TLS.KeyFile,
				CAFile:   config.TLS.CAFile,
				CertData: config.TLS.CrtData,
				KeyData:  config.TLS.KeyData,
				CAData:   config.TLS.CAData,
			},
		}
	}

	restConfig.Burst = MaxBurst
	restConfig.QPS = MaxQPS

	return restConfig, nil
}
