// +build k8srequired

package statusupdate

import (
	"context"
	"sync"
	"testing"

	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/apiextensions/pkg/clientset/versioned"
	"github.com/giantswarm/e2e-harness/pkg/harness"
	"github.com/giantswarm/microerror"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

type ResourceConfig struct {
	T *testing.T
}

type Resource struct {
	t *testing.T

	g8sClient versioned.Interface

	executionCount int
	mutex          sync.Mutex
}

func NewResource(config ResourceConfig) (*Resource, error) {
	if config.T == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.T must not be empty", config)
	}

	var g8sClient versioned.Interface
	{
		restConfig, err := clientcmd.BuildConfigFromFlags("", harness.DefaultKubeConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		g8sClient, err = versioned.NewForConfig(restConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	r := &Resource{
		t: config.T,

		g8sClient: g8sClient,

		executionCount: 0,
		mutex:          sync.Mutex{},
	}

	return r, nil
}

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	defer func() { r.executionCount++ }()

	var customResource v1alpha1.NodeConfig
	{
		curObj := obj.(*v1alpha1.NodeConfig)

		newObj, err := r.g8sClient.CoreV1alpha1().NodeConfigs(curObj.GetNamespace()).Get(curObj.GetName(), metav1.GetOptions{})
		if err != nil {
			r.t.Fatal("expected", nil, "got", err)
		}

		customResource = *newObj
	}

	if r.executionCount == 0 {
		newCondition := v1alpha1.NodeConfigStatusCondition{
			Status: conditionStatus,
			Type:   conditionType,
		}
		customResource.Status.Conditions = append(customResource.Status.Conditions, newCondition)

		_, err := r.g8sClient.CoreV1alpha1().NodeConfigs(customResource.GetNamespace()).UpdateStatus(&customResource)
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
