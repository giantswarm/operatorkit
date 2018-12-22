package controller

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	"github.com/giantswarm/operatorkit/controller/context/finalizerskeptcontext"
)

const (
	finalizerPrefix = "operatorkit.giantswarm.io"
)

type patchSpec struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

func (f *Controller) addFinalizer(ctx context.Context, obj interface{}) (bool, error) {
	// We get the accessor of the object which we got passed from the framework.
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return false, microerror.Mask(err)
	}
	// We check if the object has a finalizer here, to avoid unnecessary calls to
	// the k8s api.
	if containsFinalizer(accessor.GetFinalizers(), getFinalizerName(f.name)) {
		return false, nil // object already has the finalizer.
	}

	var stopReconciliation bool
	{
		o := func() error {
			// We get an up to date version of our object from k8s and parse the
			// response from the RESTClient to runtime object.
			obj, err := f.restClient.Get().AbsPath(accessor.GetSelfLink()).Do().Get()
			if runtime.IsNotRegisteredError(err) {
				return microerror.Mask(invalidRESTClientError)
			} else if err != nil {
				return microerror.Mask(err)
			}

			patch, stop, err := createAddFinalizerPatch(obj, f.name)
			if err != nil {
				return microerror.Mask(err)
			}
			if patch == nil {
				stopReconciliation = stop

				// When patch is empty, there nothing to do. We trust
				// createAddFinalizerPatch with the decision to stop reconciliation.
				return nil
			}

			p, err := json.Marshal(patch)
			if err != nil {
				return microerror.Mask(err)
			}
			err = f.restClient.Patch(types.JSONPatchType).AbsPath(accessor.GetSelfLink()).Body(p).Do().Error()
			if err != nil {
				return microerror.Mask(err)
			}

			stopReconciliation = true

			return nil
		}

		err = backoff.Retry(o, f.backOffFactory())
		if err != nil {
			return false, microerror.Mask(err)
		}
	}

	return stopReconciliation, nil
}

// removeFinalizer receives an object and tries to remove its finalizer which
// was set by operatorkit. The removal of a finalizer will be retried and a fresh
// object will get fetched from k8s if the ResourceVersion is out of date.
func (c *Controller) removeFinalizer(ctx context.Context, obj interface{}) error {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return microerror.Mask(err)
	}
	finalizerName := getFinalizerName(c.name)
	selfLink := accessor.GetSelfLink()

	// The control flow primitives operatorkit provides supports the mechanism of
	// keeping finalizers. This is especially useful when delete events should be
	// replayed. In case we see such a request via the dispatched context, we skip
	// the finalizer removal.
	if finalizerskeptcontext.IsKept(ctx) {
		c.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("did not remove finalizer '%s'", finalizerName))
		c.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("finalizer '%s' requested to be kept", finalizerName))

		return nil
	}

	// The reconciled object has no finalizer being set. This could have several
	// reasons. All these cases should not be harmful in general, so we ignore
	// them.
	//
	//     - We are migrating and an old object never got reconciled before
	//       deletion.
	//     - The operator wasn't running and our first interaction with the object
	//       is its deletion.
	//     - The object has another finalizer set and we removed ours already.
	//
	if !containsFinalizer(accessor.GetFinalizers(), finalizerName) {
		c.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("did not remove finalizer '%s'", finalizerName))
		c.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("finalizer '%s' not found", finalizerName))

		return nil
	}

	{
		c.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("removing finalizer '%s'", finalizerName))

		o := func() error {
			newObject, err := c.restClient.Get().AbsPath(selfLink).Do().Get()
			if errors.IsNotFound(err) {
				// The reconciled object is already gone. Nothing to do anymore.
				return nil
			} else if err != nil {
				return microerror.Mask(err)
			}

			newAccessor, err := meta.Accessor(newObject)
			if err != nil {
				return microerror.Mask(err)
			}
			patch := []patchSpec{
				{
					Op:    "replace",
					Value: removeFinalizer(newAccessor.GetFinalizers(), finalizerName),
					Path:  "/metadata/finalizers",
				},
			}

			p, err := json.Marshal(patch)
			if err != nil {
				return microerror.Mask(err)
			}
			err = c.restClient.Patch(types.JSONPatchType).AbsPath(selfLink).Body(p).Do().Error()
			if err != nil {
				return microerror.Mask(err)
			}

			return nil
		}
		b := c.backOffFactory()

		err = backoff.Retry(o, b)
		if err != nil {
			return microerror.Mask(err)
		}

		c.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("removed finalizer '%s'", finalizerName))
	}

	return nil
}

func containsFinalizer(finalizers []string, finalizer string) bool {
	for _, f := range finalizers {
		if f == finalizer {
			return true
		}
	}
	return false
}

func createAddFinalizerPatch(obj interface{}, operatorName string) (patch []patchSpec, stopReconciliation bool, err error) {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return nil, false, microerror.Mask(err)
	}
	if accessor.GetDeletionTimestamp() != nil {
		return nil, true, nil // object has been marked for deletion, we should ignore it.
	}
	finalizerName := getFinalizerName(operatorName)
	if containsFinalizer(accessor.GetFinalizers(), finalizerName) {
		return nil, false, nil // object already has the finalizer.
	}
	patch = []patchSpec{}
	if len(accessor.GetFinalizers()) == 0 {
		createPatch := patchSpec{
			Op:    "add",
			Value: []string{},
			Path:  "/metadata/finalizers",
		}
		patch = append(patch, createPatch)
	}

	addPatch := patchSpec{
		Op:    "add",
		Value: finalizerName,
		Path:  "/metadata/finalizers/-",
	}
	patch = append(patch, addPatch)

	testResourceVersionPatch := patchSpec{
		Op:    "test",
		Value: accessor.GetResourceVersion(),
		Path:  "/metadata/resourceVersion",
	}
	patch = append(patch, testResourceVersionPatch)

	return patch, true, nil
}

func getFinalizerName(name string) string {
	return fmt.Sprintf("%s/%s", finalizerPrefix, name)
}

func removeFinalizer(finalizers []string, finalizer string) []string {
	for i, f := range finalizers {
		if f == finalizer {
			finalizers = append(finalizers[:i], finalizers[i+1:]...)
			break
		}
	}

	return finalizers
}
