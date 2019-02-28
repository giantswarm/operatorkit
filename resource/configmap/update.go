package app

import (
	"context"
	"fmt"

	"github.com/giantswarm/apiextensions/pkg/apis/application/v1alpha1"
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/release-operator/service/controller/key"
)

func (r *Resource) ApplyUpdateChange(ctx context.Context, obj, updateChange interface{}) error {
	appCR, err := key.ToAppCR(updateChange)
	if err != nil {
		return microerror.Mask(err)
	}

	if appCR != nil {
		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("updating App CR %#q in namespace %#q", appCR.Name, appCR.Namespace))

		_, err = r.g8sClient.ApplicationV1alpha1().Apps(appCR.Namespace).Update(appCR)
		if err != nil {
			return microerror.Mask(err)
		}

		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("updated App CR %#q in namespace %#q", appCR.Name, appCR.Namespace))
	}

	return nil
}

func (r *Resource) newUpdateChange(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	currentAppCRs, err := toAppCRs(currentState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	desiredAppCRs, err := toAppCRs(desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var appCRsToUpdate []*v1alpha1.App
	{
		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("computing App CRs to update"))

		for _, c := range currentAppCRs {
			for _, d := range desiredAppCRs {
				m := newAppCRToUpdate(c, d)
				if m != nil {
					appCRsToUpdate = append(appCRsToUpdate, m)
				}
			}
		}

		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("computed %d App CRs to update", appCRsToUpdate))
	}

	return appCRsToUpdate, nil
}
