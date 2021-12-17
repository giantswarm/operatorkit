//go:build k8srequired
// +build k8srequired

package statusupdate

import (
	"context"
	"sync"
	"testing"

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "github.com/giantswarm/operatorkit/v6/api/v1"
	"github.com/giantswarm/operatorkit/v6/integration/env"
)

const (
	conditionStatus = "testStatus"
	conditionType   = "testType"
)

type ResourceConfig struct {
	T *testing.T
}

type Resource struct {
	t *testing.T

	ctrlClient client.Client

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
				apiextensionsv1.AddToScheme,
				v1.AddToScheme,
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

		ctrlClient: k8sClient.CtrlClient(),

		executionCount: 0,
		mutex:          sync.Mutex{},
	}

	return r, nil
}

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	defer func() { r.executionCount++ }()

	objTyped := obj.(*v1.Example)

	if r.executionCount == 0 {
		o := func() error {
			var currentObject v1.Example
			err := r.ctrlClient.Get(ctx, client.ObjectKey{
				Namespace: objTyped.Namespace,
				Name:      objTyped.Name,
			}, &currentObject)
			if err != nil {
				return microerror.Mask(err)
			}

			newCondition := v1.ExampleCondition{
				LastTransitionTime: metav1.Now(),
				Status:             conditionStatus,
				Type:               conditionType,
			}
			currentObject.Status.Conditions = append(currentObject.Status.Conditions, newCondition)

			err = r.ctrlClient.Status().Update(ctx, &currentObject)
			if err != nil {
				return microerror.Mask(err)
			}

			return nil
		}
		b := backoff.NewExponential(backoff.ShortMaxWait, backoff.ShortMaxInterval)

		err := backoff.Retry(o, b)
		if err != nil {
			r.t.Fatal("expected", nil, "got", err)
		}
	} else {
		if len(objTyped.Status.Conditions) != 1 {
			r.t.Fatalf("expected one status condition but got %d", len(objTyped.Status.Conditions))
		}
		if objTyped.Status.Conditions[0].Status != conditionStatus {
			r.t.Fatalf("expected status condition status %#q but got %#q", conditionStatus, objTyped.Status.Conditions[0].Status)
		}
		if objTyped.Status.Conditions[0].Type != conditionType {
			r.t.Fatalf("expected status condition type %#q but got %#q", conditionType, objTyped.Status.Conditions[0].Type)
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
