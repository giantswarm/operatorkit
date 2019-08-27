package secretresource

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Resource) ApplyDeleteChange(ctx context.Context, obj, deleteChange interface{}) error {
	secretsToDelete, err := toSecrets(deleteChange)
	if err != nil {
		return microerror.Mask(err)
	}

	for _, secret := range secretsToDelete {
		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("deleting Secret %#q in namespace %#q", secret.Name, secret.Namespace))

		err := r.k8sClient.CoreV1().Secrets(secret.Namespace).Delete(secret.Name, &metav1.DeleteOptions{})
		if apierrors.IsNotFound(err) {
			r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("already deleted Secret %#q in namespace %#q", secret.Name, secret.Namespace))
		} else if err != nil {
			return microerror.Mask(err)
		} else {
			r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("deleted Secret %#q in namespace %#q", secret.Name, secret.Namespace))
		}
	}

	return nil
}

func (r *Resource) newDeleteChangeForDeletePatch(ctx context.Context, obj, currentState, desiredState interface{}) ([]*corev1.Secret, error) {
	currentSecrets, err := toSecrets(currentState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return currentSecrets, nil
}

func (r *Resource) newDeleteChangeForUpdatePatch(ctx context.Context, obj, currentState, desiredState interface{}) ([]*corev1.Secret, error) {
	currentSecrets, err := toSecrets(currentState)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	desiredSecrets, err := toSecrets(desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var secretsToDelete []*corev1.Secret
	{
		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("computing Secrets to delete"))

		for _, c := range currentSecrets {
			if !containsSecret(c, desiredSecrets) {
				secretsToDelete = append(secretsToDelete, c)
			}
		}

		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("computed %d Secrets to delete", len(secretsToDelete)))
	}

	return secretsToDelete, nil
}
