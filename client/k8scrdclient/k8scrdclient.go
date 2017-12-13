package k8scrdclient

import (
	"context"

	"github.com/cenkalti/backoff"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Config struct {
	K8sExtClient apiextensionsclient.Interface
	Logger       micrologger.Logger
}

func DefaultConfig() Config {
	return Config{
		K8sExtClient: nil,
		Logger:       nil,
	}
}

type CRDClient struct {
	k8sExtClient apiextensionsclient.Interface
	logger       micrologger.Logger
}

func New(config Config) (*CRDClient, error) {
	if config.K8sExtClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.K8sExtClient must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}

	crdClient := &CRDClient{
		k8sExtClient: config.K8sExtClient,
		logger:       config.Logger,
	}

	return crdClient, nil
}

// EnsureCreated ensures the given CRD exists, is active (aka. established) and
// does not have conflicting names.
func (c *CRDClient) EnsureCreated(ctx context.Context, customResource *apiextensionsv1beta1.CustomResourceDefinition, backOff backoff.BackOff) error {
	_, err := c.k8sExtClient.ApiextensionsV1beta1().CustomResourceDefinitions().Create(customResource)
	if errors.IsAlreadyExists(err) {
		// Fall trough. We need to check CRD status.
	} else if err != nil {
		return microerror.Mask(err)
	}

	operation := func() error {
		manifest, err := c.k8sExtClient.ApiextensionsV1beta1().CustomResourceDefinitions().Get(customResource.Name, metav1.GetOptions{})
		if err != nil {
			return microerror.Mask(err)
		}

		for _, cond := range manifest.Status.Conditions {
			switch cond.Type {
			case apiextensionsv1beta1.Established:
				if cond.Status == apiextensionsv1beta1.ConditionTrue {
					return nil
				}
			case apiextensionsv1beta1.NamesAccepted:
				if cond.Status == apiextensionsv1beta1.ConditionFalse {
					return microerror.Maskf(nameConflictError, cond.Reason)
				}
			}
		}

		return microerror.Mask(notEstablishedError)
	}

	err = backoff.Retry(operation, backOff)
	if err != nil {
		deleteErr := c.k8sExtClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(customResource.Name, nil)
		if deleteErr != nil {
			return microerror.Mask(deleteErr)
		}

		return microerror.Mask(err)
	}

	return nil
}

// EnsureDeleted ensures the given CRD does not exist.
func (c *CRDClient) EnsureDeleted(ctx context.Context, customResource *apiextensionsv1beta1.CustomResourceDefinition, backOff backoff.BackOff) error {
	operation := func() error {
		err := c.k8sExtClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(customResource.Name, nil)
		if errors.IsNotFound(err) {
			// Fall trough. We reached our goal.
		} else if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}

	err := backoff.Retry(operation, backOff)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
