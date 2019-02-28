package app

import (
	"context"
	"fmt"

	"github.com/giantswarm/apiextensions/pkg/apis/application/v1alpha1"
	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// ApplyCreateChange ensures the App CR is created in the k8s api.
func (r *Resource) ApplyCreateChange(ctx context.Context, obj, createChange interface{}) error {
	appCRs, err := toAppCRs(createChange)
	if err != nil {
		return microerror.Mask(err)
	}

	for _, appCR := range appCRs {
		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("creating App CR %#q in namespace %#q", appCR.Name, appCR.Namespace))

		_, err = r.g8sClient.ApplicationV1alpha1().Apps(appCR.Namespace).Create(appCR)
		if apierrors.IsAlreadyExists(err) {
			r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("already created App CR %#q in namespace %#q", appCR.Name, appCR.Namespace))
		} else if err != nil {
			return microerror.Mask(err)
		} else {
			r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("created App CR %#q in namespace %#q", appCR.Name, appCR.Namespace))
		}
	}

	return nil
}

func (r *Resource) newCreateChange(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	currentAppCRs, err := toAppCRs(currentState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	desiredAppCRs, err := toAppCRs(desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var appCRsToCreate []*v1alpha1.App
	{
		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("computing App CRs to create "))

		for _, d := range desiredAppCRs {
			if !containsAppCR(d, currentAppCRs) {
				appCRsToCreate = append(appCRsToCreate, d)
			}
		}

		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("computed %d App CRs to create", len(appCRsToCreate)))
	}

	return appCRsToCreate, nil
}
