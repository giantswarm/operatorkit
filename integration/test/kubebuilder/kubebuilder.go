package kubebuilder

import (
	"testing"
	"time"

	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/micrologger/microloggertest"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	operatorkitcontroller "github.com/giantswarm/operatorkit/v4/pkg/controller"
)

func Test_Kubebuilder_Integration_Basic(t *testing.T) {
	logger := microloggertest.New()
	scheme := runtime.NewScheme()
	syncPeriod := time.Second * 5
	options := controllerruntime.Options{
		Scheme:             scheme,
		SyncPeriod:         &syncPeriod,
		RetryPeriod:        nil,
		MetricsBindAddress: "0",
	}
	mgr, err := controllerruntime.NewManager(controllerruntime.GetConfigOrDie(), options)
	if err != nil {
		t.Fatal(err)
	}
	reconciler := &CronJobReconciler{
		Log: logger,
	}

	controllerOptions := controller.Options{
		MaxConcurrentReconciles: 1,
		Reconciler:              reconciler,
	}

	err = reconciler.SetupWithManager(mgr, controllerOptions)
	if err != nil {
		t.Fatal(err)
	}

	_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: types.NamespacedName{
		Namespace: "default",
		Name:      "some",
	}})
	if err != nil {
		t.Fatal(err)
	}
}

// CronJobReconciler reconciles a CronJob object
type CronJobReconciler struct {
	operatorkitcontroller.Controller
	Log micrologger.Logger
}
