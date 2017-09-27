package crd

import (
	"time"

	"github.com/cenk/backoff"
	"github.com/giantswarm/microerror"
	apiv1 "k8s.io/api/core/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

const (
	// ResyncPeriod is the interval at which the Informer cache is invalidated,
	// and the lister function is called.
	ResyncPeriod = 1 * time.Minute
)

// InformerConfig represents the configuration used to create a new Informer.
type InformerConfig struct {
	// Dependencies.
	BackOff   backoff.BackOff
	CRDClient apiextensionsclient.Interface

	// Settings.
	ResyncPeriod time.Duration
}

// DefaultInformerConfig provides a default configuration to create a new
// Informer by best effort.
func DefaultInformerConfig() InformerConfig {
	return InformerConfig{
		// Dependencies.
		BackOff:   nil,
		CRDClient: nil,

		// Settings.
		ResyncPeriod: ResyncPeriod,
	}
}

type Informer struct {
	// Dependencies.
	backOff   backoff.BackOff
	crdClient apiextensionsclient.Interface

	// Settings.
	resyncPeriod time.Duration
}

// New creates a new Informer.
func NewInformer(config InformerConfig) (*Informer, error) {
	// Dependencies.
	if config.BackOff == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.BackOff must not be empty")
	}
	if config.CRDClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.CRDClient must not be empty")
	}

	// Settings.
	if config.ResyncPeriod == 0 {
		return nil, microerror.Maskf(invalidConfigError, "config.ResyncPeriod must not be empty")
	}

	newInformer := &Informer{
		// Settings.
		backOff:   config.BackOff,
		crdClient: config.CRDClient,

		// Settings.
		resyncPeriod: config.ResyncPeriod,
	}

	return newInformer, nil
}

func (i *Informer) CreateCRD(CRD *CRD) error {
	_, err := i.crdClient.ApiextensionsV1beta1().CustomResourceDefinitions().Create(CRD.NewResource())
	if err != nil {
		return microerror.Mask(err)
	}

	operation := func() error {
		manifest, err := i.crdClient.ApiextensionsV1beta1().CustomResourceDefinitions().Get(CRD.Name(), metav1.GetOptions{})
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

	err = backoff.Retry(operation, i.backOff)
	if err != nil {
		deleteErr := i.crdClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(CRD.Name(), nil)
		if deleteErr != nil {
			return microerror.Mask(deleteErr)
		}

		return microerror.Mask(err)
	}

	return nil
}

func (i *Informer) NewController(CRD *CRD, resourceEventHandler cache.ResourceEventHandler, zeroObjectFactory ZeroObjectFactory) cache.Controller {
	listWatch := cache.NewListWatchFromClient(
		i.crdClient.ApiextensionsV1beta1().RESTClient(),
		CRD.Plural(),
		apiv1.NamespaceAll,
		fields.Everything(),
	)

	_, controller := cache.NewInformer(listWatch, zeroObjectFactory.NewObject(), i.resyncPeriod, resourceEventHandler)

	return controller
}
