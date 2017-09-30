package crd

import (
	"context"

	"github.com/cenk/backoff"
	"github.com/giantswarm/microerror"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Ensure(ctx context.Context, CRD *CRD, crdClient apiextensionsclient.Interface, backOff backoff.BackOff) error {
	_, err := crdClient.ApiextensionsV1beta1().CustomResourceDefinitions().Create(CRD.NewResource())
	if errors.IsAlreadyExists(err) {
		return nil
	} else if err != nil {
		return microerror.Mask(err)
	}

	operation := func() error {
		manifest, err := crdClient.ApiextensionsV1beta1().CustomResourceDefinitions().Get(CRD.Name(), metav1.GetOptions{})
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
		deleteErr := crdClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(CRD.Name(), nil)
		if deleteErr != nil {
			return microerror.Mask(deleteErr)
		}

		return microerror.Mask(err)
	}

	return nil
}
