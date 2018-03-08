package framework

import (
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
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
		return false, err
	}
	if containsFinalizer(accessor.GetFinalizers(), finalizerName) {
		return false, nil // resource already has the finalizer.
	}

	c, err := rest.RESTClientFor(&rest.Config{})
	if err != nil {
		return false, err
	}
	patch := fmt.Sprintf(`{"op":"add","value":%q,"path":"/metadata/finalizers/1"}`, finalizerName)
	res := c.Patch(types.JSONPatchType).
		AbsPath(accessor.GetSelfLink()).
		Body(patch).
		Do()
	if errors.IsForbidden(res.Error()) {
		return true, nil // we tried to add the finalizer again after removing it.
	} else if res.Error() != nil {
		return false, res.Error()
	}

	return true, nil
}

func (f *Framework) removeFinalizer(obj interface{}) error {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return err
	}
	if !containsFinalizer(accessor.GetFinalizers(), finalizerName) {
		return nil // resource has no finalizer, probably migration.
	}
	c, err := rest.RESTClientFor(&rest.Config{})
	if err != nil {
		return err
	}

	patch := patchSpec{
		Op:    "replace",
		Value: removeFinalizer(accessor.GetFinalizers(), finalizerName),
		Path:  "/metadata/finalizers",
	}
	p, err := json.Marshal(patch)
	if err != nil {
		return err
	}
	res := c.Patch(types.JSONPatchType).
		AbsPath(accessor.GetSelfLink()).
		Body(p).
		Do()
	if res.Error() != nil {
		return res.Error()
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
