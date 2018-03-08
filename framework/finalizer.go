package framework

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/giantswarm/microerror"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
)

const (
	finalizerName = "operatorkit.giantswarm.io"
)

type patchSpec struct {
	Op    string   `json:"op"`
	Path  string   `json:"path"`
	Value []string `json:"value"`
}

func (f *Framework) addFinalizer(obj interface{}) (bool, error) {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return false, microerror.Mask(err)
	}
	if accessor.GetDeletionTimestamp() != nil {
		return true, nil // object has been marked for deletion, we should ignore it.
	}
	if containsFinalizer(accessor.GetFinalizers(), finalizerName) {
		return false, nil // object already has the finalizer.
	}

	c, err := rest.RESTClientFor(&rest.Config{})
	if err != nil {
		return false, microerror.Mask(err)
	}
	patch := fmt.Sprintf(`{"op":"add","value":%q,"path":"/metadata/finalizers/1"}`, finalizerName)
	res := c.Patch(types.JSONPatchType).AbsPath(accessor.GetSelfLink()).Body(patch).Do()
	if res.Error() != nil {
		return false, microerror.Mask(res.Error())
	}

	return true, nil
}

func (f *Framework) removeFinalizer(ctx context.Context, obj interface{}) error {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return microerror.Mask(err)
	}
	if !containsFinalizer(accessor.GetFinalizers(), finalizerName) {
		f.logger.LogCtx(ctx, "function", "removeFinalizer", "level", "warning", "message", "object is missing a finalizer")
		return nil // object has no finalizer, probably migration.
	}
	c, err := rest.RESTClientFor(&rest.Config{})
	if err != nil {
		return microerror.Mask(err)
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
	res := c.Patch(types.JSONPatchType).AbsPath(accessor.GetSelfLink()).Body(p).Do()
	if res.Error() != nil {
		return microerror.Mask(res.Error())
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

func removeFinalizer(finalizers []string, finalizer string) []string {
	for i, f := range finalizers {
		if f == finalizer {
			finalizers = append(finalizers[:i], finalizers[i+1:]...)
			break
		}
	}
	return finalizers
}
