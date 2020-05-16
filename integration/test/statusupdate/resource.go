// +build k8srequired

package statusupdate

import (
	"context"
	"sync"
	"testing"

	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/k8sclient/v3/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/operatorkit/integration/env"
)

type ResourceConfig struct {
	T *testing.T
}

type Resource struct {
	t *testing.T

	k8sClient k8sclient.Interface

	executionCount int
	mutex          sync.Mutex
}

func NewResource(config ResourceConfig) (*Resource, error) {
	if config.T == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.T must not be empty", config)
	}

	var err error

	var newLogger micrologger.Logger
	{
		c := micrologger.Config{}

		newLogger, err = micrologger.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var k8sClient *k8sclient.Clients
	{
		c := k8sclient.ClientsConfig{
			SchemeBuilder: k8sclient.SchemeBuilder{
				v1alpha1.AddToScheme,
			},
			Logger: newLogger,

			KubeConfigPath: env.KubeConfigPath(),
		}

		k8sClient, err = k8sclient.NewClients(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	r := &Resource{
		t: config.T,

		k8sClient: k8sClient,

		executionCount: 0,
		mutex:          sync.Mutex{},
	}

	return r, nil
}

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	defer func() { r.executionCount++ }()

	var customResource v1alpha1.DrainerConfig
	{
		curObj := obj.(*v1alpha1.DrainerConfig)

		newObj, err := r.k8sClient.G8sClient().CoreV1alpha1().DrainerConfigs(curObj.GetNamespace()).Get(curObj.GetName(), metav1.GetOptions{})
		if err != nil {
			r.t.Fatal("expected", nil, "got", err)
		}

		customResource = *newObj
	}

	if r.executionCount == 0 {
		newCondition := v1alpha1.DrainerConfigStatusCondition{
			LastHeartbeatTime:  metav1.Now(),
			LastTransitionTime: metav1.Now(),
			Status:             conditionStatus,
			Type:               conditionType,
		}
		customResource.Status.Conditions = append(customResource.Status.Conditions, newCondition)

		_, err := r.k8sClient.G8sClient().CoreV1alpha1().DrainerConfigs(customResource.GetNamespace()).UpdateStatus(&customResource)
		if err != nil {
			r.t.Fatal("expected", nil, "got", err)
		}
	} else {
		if len(customResource.Status.Conditions) != 1 {
			r.t.Fatalf("expected one status condition but got %d", len(customResource.Status.Conditions))
		}
		if customResource.Status.Conditions[0].Status != conditionStatus {
			r.t.Fatalf("expected status condition status %#q but got %#q", conditionStatus, customResource.Status.Conditions[0].Status)
		}
		if customResource.Status.Conditions[0].Type != conditionType {
			r.t.Fatalf("expected status condition type %#q but got %#q", conditionType, customResource.Status.Conditions[0].Type)
		}
	}

	return nil
}

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	return nil
}

func (r *Resource) Name() string {
	return "statusupdate"
}
