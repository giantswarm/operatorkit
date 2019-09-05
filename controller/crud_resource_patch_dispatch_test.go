// NOTE the CRUD resource has moved to operatorkit/resource/crud. The code below
// is DEPRECATED and only kept for backward compatibility.
package controller

import (
	"context"
	"fmt"
	"testing"

	"github.com/giantswarm/micrologger/microloggertest"
)

func Test_CRUDResource_PatchDispatch(t *testing.T) {
	testCases := []struct {
		Ops CRUDResourceOps
	}{
		// Case 0
		{
			Ops: &testCRUDResourceOpsPatchDispatch{
				Patch: func() *Patch {
					return nil
				}(),
			},
		},

		// Case 1
		{
			Ops: &testCRUDResourceOpsPatchDispatch{
				Patch: func() *Patch {
					p := NewPatch()

					p.SetCreateChange(true)

					return p
				}(),
			},
		},

		// Case 2
		{
			Ops: &testCRUDResourceOpsPatchDispatch{
				Patch: func() *Patch {
					p := NewPatch()

					p.SetDeleteChange(true)

					return p
				}(),
			},
		},

		// Case 3
		{
			Ops: &testCRUDResourceOpsPatchDispatch{
				Patch: func() *Patch {
					p := NewPatch()

					p.SetDeleteChange(true)

					return p
				}(),
			},
		},

		// Case 4
		{
			Ops: &testCRUDResourceOpsPatchDispatch{
				Patch: func() *Patch {
					p := NewPatch()

					p.SetDeleteChange(true)
					p.SetUpdateChange(true)

					return p
				}(),
			},
		},

		// Case 5
		{
			Ops: &testCRUDResourceOpsPatchDispatch{
				Patch: func() *Patch {
					p := NewPatch()

					p.SetCreateChange(true)
					p.SetDeleteChange(true)
					p.SetUpdateChange(true)

					return p
				}(),
			},
		},
	}

	for i, tc := range testCases {
		c := CRUDResourceConfig{
			Logger: microloggertest.New(),
			Ops:    tc.Ops,
		}
		r, err := NewCRUDResource(c)
		if err != nil {
			t.Fatalf("test %d: unexpected NewCRUDResource error = %#v", i, err)
		}

		err = r.EnsureCreated(context.Background(), nil)
		if err != nil {
			t.Fatalf("test %d: unexpected EnsureCreated error = %#v", i, err)
		}

		err = r.EnsureDeleted(context.Background(), nil)
		if err != nil {
			t.Fatalf("test %d: unexpected EnsureDeleted error = %#v", i, err)
		}
	}
}

type testCRUDResourceOpsPatchDispatch struct {
	Patch *Patch
}

func (r *testCRUDResourceOpsPatchDispatch) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	return nil, nil
}

func (r *testCRUDResourceOpsPatchDispatch) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	return nil, nil
}

func (r *testCRUDResourceOpsPatchDispatch) NewUpdatePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*Patch, error) {
	return r.Patch, nil
}

func (r *testCRUDResourceOpsPatchDispatch) NewDeletePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*Patch, error) {
	return r.Patch, nil
}

func (r *testCRUDResourceOpsPatchDispatch) Name() string {
	return "testCRUDResourceOpsPatchDispatch"
}

func (r *testCRUDResourceOpsPatchDispatch) ApplyCreateChange(ctx context.Context, obj, createState interface{}) error {
	createChange, ok := r.Patch.getCreateChange()
	if ok {
		if createChange != createState {
			panic(fmt.Sprintf("expected '%s' got '%s'", createChange, createState))
		}
	}

	return nil
}

func (r *testCRUDResourceOpsPatchDispatch) ApplyDeleteChange(ctx context.Context, obj, deleteState interface{}) error {
	deleteChange, ok := r.Patch.getDeleteChange()
	if ok {
		if deleteChange != deleteState {
			panic(fmt.Sprintf("expected '%s' got '%s'", deleteChange, deleteState))
		}
	}

	return nil
}

func (r *testCRUDResourceOpsPatchDispatch) ApplyUpdateChange(ctx context.Context, obj, updateState interface{}) error {
	updateChange, ok := r.Patch.getUpdateChange()
	if ok {
		if updateChange != updateState {
			panic(fmt.Sprintf("expected '%s' got '%s'", updateChange, updateState))
		}
	}

	return nil
}
