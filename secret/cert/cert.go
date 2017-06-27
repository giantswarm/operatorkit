package cert

import (
	"fmt"
	"time"

	"github.com/giantswarm/certificatetpr"
	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
)

const (
	// WatchTimeOut is the time to wait on watches against the Kubernetes API
	// before giving up and throwing an error.
	WatchTimeOut = 90 * time.Second
)

// Config represents the configuration used to create a service.
type Config struct {
	// Dependencies.
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger
}

// DefaultConfig provides a default configuration to create a new service by
// best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		K8sClient: nil,
		Logger:    nil,
	}
}

// New creates a new configured service.
func New(config Config) (*Service, error) {
	// Dependencies.
	if config.K8sClient == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.K8sClient must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.Logger must not be empty")
	}

	newService := &Service{
		// Dependencies.
		k8sClient: config.K8sClient,
		logger:    config.Logger,
	}

	return newService, nil
}

// Service implements the version service interface.
type Service struct {
	// Dependencies.
	k8sClient kubernetes.Interface
	logger    micrologger.Logger
}

// SearchCerts watches for all secrets of a cluster  and returns it as
// assets bundle.
func (s *Service) SearchCerts(clusterID string) (certificatetpr.AssetsBundle, error) {
	assetsBundle := make(certificatetpr.AssetsBundle)

	for _, componentName := range certificatetpr.ClusterComponents {
		ab, err := s.SearchCertsForComponent(clusterID, componentName.String())
		if err != nil {
			return nil, microerror.MaskAny(err)
		}

		for k, v := range ab {
			assetsBundle[k] = v
		}
	}

	return assetsBundle, nil
}

// SearchCertsForComponent watches for secrets of a single cluster component and
// returns it as assets bundle.
func (s *Service) SearchCertsForComponent(clusterID, componentName string) (certificatetpr.AssetsBundle, error) {
	// TODO we should also do a list. In case the secrets have already been
	// created we might miss them with only watching.
	watcher, err := s.k8sClient.Core().Secrets(api.NamespaceDefault).Watch(apismetav1.ListOptions{
		// Select only secrets that match the given component and the given cluster
		// clusterID.
		LabelSelector: fmt.Sprintf(
			"%s=%s, %s=%s",
			certificatetpr.ComponentLabel,
			componentName,
			certificatetpr.ClusterIDLabel,
			clusterID,
		),
	})
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	assetsBundle := make(certificatetpr.AssetsBundle)

	defer watcher.Stop()
	for {
		select {
		case event, ok := <-watcher.ResultChan():
			if !ok {
				return nil, microerror.MaskAnyf(secretsRetrievalFailedError, "secrets channel was already closed")
			}

			switch event.Type {
			case watch.Added:
				secret := event.Object.(*v1.Secret)
				component := certificatetpr.ClusterComponent(secret.Labels[certificatetpr.ComponentLabel])

				if !certificatetpr.ValidComponent(component, certificatetpr.ClusterComponents) {
					return nil, microerror.MaskAnyf(secretsRetrievalFailedError, "unknown clusterComponent %s", component)
				}

				for _, assetType := range certificatetpr.TLSAssetTypes {
					asset, ok := secret.Data[assetType.String()]
					if !ok {
						return nil, microerror.MaskAnyf(secretsRetrievalFailedError, "malformed secret was missing %v asset", assetType)
					}

					assetsBundle[certificatetpr.AssetsBundleKey{component, assetType}] = asset
				}

				return assetsBundle, nil
			case watch.Deleted:
				// Noop. Ignore deleted events. These are handled by the certificate
				// operator.
			case watch.Error:
				return nil, microerror.MaskAnyf(secretsRetrievalFailedError, "there was an error in the watcher: %v", apierrors.FromObject(event.Object))
			}
		case <-time.After(WatchTimeOut):
			return nil, microerror.MaskAnyf(secretsRetrievalFailedError, "timed out waiting for secrets")
		}
	}
}
