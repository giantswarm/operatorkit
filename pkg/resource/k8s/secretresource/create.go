package secretresource

import (
	"context"

	"github.com/giantswarm/microerror"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ApplyCreateChange ensures the Secret is created in the k8s api.
func (r *Resource) ApplyCreateChange(ctx context.Context, obj, createChange interface{}) error {
	secrets, err := toSecrets(createChange)
	if err != nil {
		return microerror.Mask(err)
	}

	for _, secret := range secrets {
		r.logger.Debugf(ctx, "creating Secret %#q in namespace %#q", secret.Name, secret.Namespace)

		_, err = r.k8sClient.CoreV1().Secrets(secret.Namespace).Create(ctx, secret, metav1.CreateOptions{})
		if apierrors.IsAlreadyExists(err) {
			r.logger.Debugf(ctx, "already created Secret %#q in namespace %#q", secret.Name, secret.Namespace)
		} else if err != nil {
			return microerror.Mask(err)
		} else {
			r.logger.Debugf(ctx, "created Secret %#q in namespace %#q", secret.Name, secret.Namespace)
		}
	}

	return nil
}

func (r *Resource) newCreateChange(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	currentSecrets, err := toSecrets(currentState)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	desiredSecrets, err := toSecrets(desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var secretsToCreate []*corev1.Secret
	{
		r.logger.Debugf(ctx, "computing Secrets to create ")

		for _, d := range desiredSecrets {
			if !containsSecret(d, currentSecrets) {
				secretsToCreate = append(secretsToCreate, d)
			}
		}

		r.logger.Debugf(ctx, "computed %d Secrets to create", len(secretsToCreate))
	}

	return secretsToCreate, nil
}
