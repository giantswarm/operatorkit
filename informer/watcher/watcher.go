package watcher

import (
	"github.com/giantswarm/microerror"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Stops watching. Will close the channel returned by ResultChan(). Releases
// any resources used by the watch.
func (w *Watcher) Stop() {
	// TODO
}

// Returns a chan which will receive all the events. If an error occurs
// or Stop() is called, this channel will be closed, in which case the
// watch should be completely cleaned up.
func (w *Watcher) ResultChan() <-chan Event {
	var mgr manager.Manager
	{
		mgrOpts := manager.Options{
			SyncPeriod: durationP(resyncPeriod),
		}

		mgr, err = manager.New(restConfig, mgrOpts)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	var ctrl controller.Controller
	{
		c := controller.Options{
			Reconciler: &reconciler{},
		}

		ctrl, err = controller.New("pod-controller", mgr, c)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	{
		err = ctrl.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForObject{})
		if err != nil {
			return microerror.Mask(err)
		}
	}

	{
		err = mgr.Start(signals.SetupSignalHandler())
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}
