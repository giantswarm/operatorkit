package secretresource

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	v1 "k8s.io/api/core/v1"
)

func (r *Resource) ApplyUpdateChange(ctx context.Context, obj, updateChange interface{}) error {
	secrets, err := toSecrets(updateChange)
	if err != nil {
		return microerror.Mask(err)
	}

	for _, secret := range secrets {
		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("updating Secret %#q in namespace %#q", secret.Name, secret.Namespace))

		_, err = r.k8sClient.CoreV1().Secrets(secret.Namespace).Update(secret)
		if err != nil {
			return microerror.Mask(err)
		}

		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("updated Secret %#q in namespace %#q", secret.Name, secret.Namespace))
	}

	return nil
}

func (r *Resource) newUpdateChange(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	currentSecrets, err := toSecrets(currentState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	desiredSecrets, err := toSecrets(desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var secretsToUpdate []*v1.Secret
	{
		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("computing Secrets to update"))

		for _, c := range currentSecrets {
			for _, d := range desiredSecrets {
				m := newSecretToUpdate(c, d)
				if m != nil {
					secretsToUpdate = append(secretsToUpdate, m)
				}
			}
		}

		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("computed %d Secrets to update", len(secretsToUpdate)))
	}

	return secretsToUpdate, nil
}
