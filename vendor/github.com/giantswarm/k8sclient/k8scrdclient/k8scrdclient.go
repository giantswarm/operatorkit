package k8scrdclient

import (
	"context"

	"github.com/giantswarm/backoff"
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

type CRDClient struct {
	k8sExtClient apiextensionsclient.Interface
	logger       micrologger.Logger
}

func New(config Config) (*CRDClient, error) {
	if config.K8sExtClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sExtClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	crdClient := &CRDClient{
		k8sExtClient: config.K8sExtClient,
		logger:       config.Logger,
	}

	return crdClient, nil
}

// EnsureCreated ensures the given CRD exists, is active (aka. established) and
// does not have conflicting names.
func (c *CRDClient) EnsureCreated(ctx context.Context, crd *apiextensionsv1beta1.CustomResourceDefinition, b backoff.Interface) error {
	var err error

	err = c.ensureCreated(ctx, crd, b)
	if err != nil {
		return microerror.Mask(err)
	}

	err = c.ensureStatusSubresourceCreated(ctx, crd, b)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

// EnsureDeleted ensures the given CRD does not exist.
func (c *CRDClient) EnsureDeleted(ctx context.Context, crd *apiextensionsv1beta1.CustomResourceDefinition, b backoff.Interface) error {
	o := func() error {
		err := c.k8sExtClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(crd.Name, nil)
		if errors.IsNotFound(err) {
			// Fall trough. We reached our goal.
		} else if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}

	err := backoff.Retry(o, b)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (c *CRDClient) ensureCreated(ctx context.Context, crd *apiextensionsv1beta1.CustomResourceDefinition, b backoff.Interface) error {
	_, err := c.k8sExtClient.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
	if errors.IsAlreadyExists(err) {
		// Fall trough. We need to check CRD status.
	} else if err != nil {
		return microerror.Mask(err)
	}

	o := func() error {
		manifest, err := c.k8sExtClient.ApiextensionsV1beta1().CustomResourceDefinitions().Get(crd.Name, metav1.GetOptions{})
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

	err = backoff.Retry(o, b)
	if err != nil {
		deleteErr := c.k8sExtClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(crd.Name, nil)
		if deleteErr != nil {
			return microerror.Mask(deleteErr)
		}

		return microerror.Mask(err)
	}

	return nil
}

// ensureStatusSubresourceCreated ensures if the CRD has a status subresource
// it is created. This is needed if a previous version of the CRD without the
// status subresource is present.
func (c *CRDClient) ensureStatusSubresourceCreated(ctx context.Context, crd *apiextensionsv1beta1.CustomResourceDefinition, b backoff.Interface) error {
	if crd.Spec.Subresources == nil || crd.Spec.Subresources.Status == nil {
		// Nothing to do.
		return nil
	}

	o := func() error {
		manifest, err := c.k8sExtClient.ApiextensionsV1beta1().CustomResourceDefinitions().Get(crd.Name, metav1.GetOptions{})
		if err != nil {
			return microerror.Mask(err)
		}

		if manifest.Spec.Subresources == nil || manifest.Spec.Subresources.Status == nil {
			crd.SetResourceVersion(manifest.ResourceVersion)
			_, err = c.k8sExtClient.ApiextensionsV1beta1().CustomResourceDefinitions().Update(crd)
			if err != nil {
				return microerror.Mask(err)
			}
		}

		return nil
	}

	err := backoff.Retry(o, b)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
