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

func (f *Framework) addFinalizer(obj interface{}) (bool, error) {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return false, microerror.Mask(err)
	}
	if accessor.GetDeletionTimestamp() != nil {
		return true, nil // object has been marked for deletion, we should ignore it.
	}
	finalizerName := getFinalizerName(f.name)
	if containsFinalizer(accessor.GetFinalizers(), finalizerName) {
		return false, nil // object already has the finalizer.
	}

	patch := patchSpec{
		Op:    "add",
		Value: finalizerName,
		Path:  "/metadata/finalizers/-",
	}
	p, err := json.Marshal(patch)
	if err != nil {
		return false, microerror.Mask(err)
	}
	operation := func() error {
		res := f.restClient.Patch(types.JSONPatchType).AbsPath(accessor.GetSelfLink()).Body(p).Do()
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

func (f *Framework) removeFinalizer(ctx context.Context, obj interface{}) error {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return microerror.Mask(err)
	}
	finalizerName := getFinalizerName(f.name)
	if !containsFinalizer(accessor.GetFinalizers(), finalizerName) {
		f.logger.LogCtx(ctx, "function", "removeFinalizer", "level", "warning", "message", "object is missing a finalizer")
		return nil // object has no finalizer, probably migration.
	}

	patch := patchSpec{
		Op:    "replace",
		Value: removeFinalizer(accessor.GetFinalizers(), finalizerName),
		Path:  "/metadata/finalizers",
	}
	p, err := json.Marshal(patch)
	if err != nil {
		return microerror.Mask(err)
	}
	operation := func() error {
		res := f.restClient.Patch(types.JSONPatchType).AbsPath(accessor.GetSelfLink()).Body(p).Do()
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
