package controller

import (
	"encoding/json"
	"fmt"

	"github.com/cenkalti/backoff"
	"github.com/giantswarm/microerror"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/types"
)

const (
	finalizerPrefix = "operatorkit.giantswarm.io"
)

type patchSpec struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

func (f *Controller) addFinalizer(obj interface{}) (bool, error) {
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
			if err != nil {
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
func (f *Controller) removeFinalizer(obj interface{}) error {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return microerror.Mask(err)
	}
	if !containsFinalizer(accessor.GetFinalizers(), getFinalizerName(f.name)) {
		// object has no finalizer set, this could have two reasons:
		// 1. We are migrating and an old object never got reconciled before deletion.
		// 2. The operator wasn't running and our first interaction with the object
		// is its deletion.
		// 3. The object has another finalizer set and we removed ours already.
		// All cases should not be harmful in general, so we ignore it.
		return nil
	}

	path := accessor.GetSelfLink()

	o := func() error {
		// We get an up to date version of our object from k8s and parse the
		// response from the RESTClient to runtime object.
		obj, err := f.restClient.Get().AbsPath(path).Do().Get()
		if errors.IsNotFound(err) {
			return nil // the object is already gone, nothing to do.
		} else if err != nil {
			return microerror.Mask(err)
		}

		patch, err := createRemoveFinalizerPatch(obj, f.name)
		if err != nil {
			return microerror.Mask(err)
		}
		if patch == nil {
			return nil
		}

		p, err := json.Marshal(patch)
		if err != nil {
			return microerror.Mask(err)
		}
		err = f.restClient.Patch(types.JSONPatchType).AbsPath(path).Body(p).Do().Error()
		if err != nil {
			return microerror.Mask(err)
		}
		return nil
	}
	err = backoff.Retry(o, f.backOffFactory())
	if err != nil {
		return microerror.Mask(err)
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

func createRemoveFinalizerPatch(obj interface{}, operatorName string) ([]patchSpec, error) {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	finalizerName := getFinalizerName(operatorName)
	if !containsFinalizer(accessor.GetFinalizers(), finalizerName) {
		return nil, nil
	}

	patch := []patchSpec{
		{
			Op:    "replace",
			Value: removeFinalizer(accessor.GetFinalizers(), finalizerName),
			Path:  "/metadata/finalizers",
		},
	}

	return patch, nil
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
