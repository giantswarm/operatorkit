package k8sclient

import (
	"context"
	"fmt"

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	v1 "k8s.io/api/core/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/k8sclient/k8scrdclient"
)

type SetupConfig struct {
	Clients *Clients
	Logger  micrologger.Logger
}

type Setup struct {
	clients *Clients
	logger  micrologger.Logger
}

func NewSetup(config SetupConfig) (*Setup, error) {
	if config.Clients == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Clients must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	s := &Setup{
		clients: config.Clients,
		logger:  config.Logger,
	}

	return s, nil
}

func (s *Setup) EnsureCRDCreated(ctx context.Context, crd *apiextensionsv1beta1.CustomResourceDefinition) error {
	var err error

	var crdClient *k8scrdclient.CRDClient
	{
		c := k8scrdclient.Config{
			K8sExtClient: s.clients.ExtClient(),
			Logger:       s.logger,
		}

		crdClient, err = k8scrdclient.New(c)
	}

	b := backoff.NewExponential(backoff.ShortMaxWait, backoff.ShortMaxInterval)

	err = crdClient.EnsureCreated(ctx, crd, b)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (s *Setup) EnsureCRDDeleted(ctx context.Context, crd *apiextensionsv1beta1.CustomResourceDefinition) error {
	var err error

	var crdClient *k8scrdclient.CRDClient
	{
		c := k8scrdclient.Config{
			K8sExtClient: s.clients.ExtClient(),
			Logger:       s.logger,
		}

		crdClient, err = k8scrdclient.New(c)
	}

	b := backoff.NewExponential(backoff.ShortMaxWait, backoff.ShortMaxInterval)

	err = crdClient.EnsureDeleted(ctx, crd, b)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (s *Setup) EnsureNamespaceCreated(ctx context.Context, namespace string) error {
	s.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("ensuring Kubernetes Namespace %#q", namespace))

	o := func() error {
		{
			n := &v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			_, err := s.clients.K8sClient().CoreV1().Namespaces().Create(n)
			if errors.IsAlreadyExists(err) {
				// fall through
			} else if err != nil {
				return microerror.Mask(err)
			}
		}

		{
			n, err := s.clients.K8sClient().CoreV1().Namespaces().Get(namespace, metav1.GetOptions{})
			if err != nil {
				return microerror.Mask(err)
			}
			if n.Status.Phase != v1.NamespaceActive {
				return microerror.Maskf(unexpectedStatusPhaseError, string(n.Status.Phase))
			}
		}

		return nil
	}
	b := backoff.NewExponential(backoff.ShortMaxWait, backoff.ShortMaxInterval)

	err := backoff.Retry(o, b)
	if err != nil {
		return microerror.Mask(err)
	}

	s.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("ensured Kubernetes Namespace %#q", namespace))

	return nil
}

func (s *Setup) EnsureNamespaceDeleted(ctx context.Context, namespace string) error {
	s.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("ensuring deletion of Kubernetes Namespace %#q", namespace))

	o := func() error {
		{
			err := s.clients.K8sClient().CoreV1().Namespaces().Delete(namespace, &metav1.DeleteOptions{})
			if errors.IsNotFound(err) {
				// fall through
			} else if err != nil {
				return microerror.Mask(err)
			}
		}

		return nil
	}
	b := backoff.NewExponential(backoff.ShortMaxWait, backoff.ShortMaxInterval)

	err := backoff.Retry(o, b)
	if err != nil {
		return microerror.Mask(err)
	}

	s.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("ensured deletion of Kubernetes Namespace %#q", namespace))

	return nil
}
