package framework

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cenkalti/backoff"
	"github.com/giantswarm/microerror"
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

func (f *Framework) addFinalizer(obj interface{}) (stopReconciliation bool, err error) {
	patch, path, result, err := createAddFinalizerPatch(obj, f.name)
	if err != nil {
		return false, microerror.Mask(err)
	}
	if patch == nil {
		return result, err
	}
	p, err := json.Marshal(patch)
	if err != nil {
		return false, microerror.Mask(err)
	}
	operation := func() error {
		res := f.restClient.Patch(types.JSONPatchType).AbsPath(path).Body(p).Do()
		if res.Error() != nil {
			return microerror.Mask(res.Error())
		}
		return nil
	}
	err = backoff.Retry(operation, f.backOffFactory())
	if err != nil {
		return false, microerror.Mask(err)
	}

	return true, nil
}

func createAddFinalizerPatch(obj interface{}, operatorName string) (patch []patchSpec, path string, stopReconciliation bool, err error) {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return nil, "", false, microerror.Mask(err)
	}
	if accessor.GetDeletionTimestamp() != nil {
		return nil, "", true, nil // object has been marked for deletion, we should ignore it.
	}
	finalizerName := getFinalizerName(operatorName)
	if containsFinalizer(accessor.GetFinalizers(), finalizerName) {
		return nil, "", false, nil // object already has the finalizer.
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
	return patch, accessor.GetSelfLink(), true, nil
}

func createRemoveFinalizerPatch(obj interface{}, operatorName string) (patch []patchSpec, path string, err error) {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return nil, "", microerror.Mask(err)
	}
	finalizerName := getFinalizerName(operatorName)
	if !containsFinalizer(accessor.GetFinalizers(), finalizerName) {
		return nil, "", nil // object has no finalizer, probably migration.
	}
	patch = []patchSpec{}
	deletePatch := patchSpec{
		Op:    "replace",
		Value: removeFinalizer(accessor.GetFinalizers(), finalizerName),
		Path:  "/metadata/finalizers",
	}
	patch = append(patch, deletePatch)
	return patch, accessor.GetSelfLink(), nil
}

func (f *Framework) removeFinalizer(ctx context.Context, obj interface{}) error {
	patch, path, err := createRemoveFinalizerPatch(obj, f.name)
	if err != nil {
		return microerror.Mask(err)
	}
	if patch == nil {
		f.logger.LogCtx(ctx, "function", "removeFinalizer", "level", "warning", "message", "object is missing a finalizer")
		return nil
	}
	p, err := json.Marshal(patch)
	if err != nil {
		return microerror.Mask(err)
	}
	operation := func() error {
		res := f.restClient.Patch(types.JSONPatchType).AbsPath(path).Body(p).Do()
		if res.Error() != nil {
			return microerror.Mask(res.Error())
		}
		return nil
	}
	err = backoff.Retry(operation, f.backOffFactory())
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
