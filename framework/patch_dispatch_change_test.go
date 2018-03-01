package framework

import (
	"context"
	"fmt"
	"testing"
)

func Test_Framework_Resource_PatchDispatchChange(t *testing.T) {
	testCases := []struct {
		ProcessMethod func(ctx context.Context, obj interface{}, rs []Resource) error
		Resources     []Resource
		ErrorMatcher  func(err error) bool
	}{
		// Case 0
		{
			ProcessMethod: ProcessDelete,
			Resources: []Resource{
				&testResourcePatchDispatchChange{
					Patch: func() *Patch {
						return nil
					}(),
				},
			},
			ErrorMatcher: nil,
		},

		// Case 1
		{
			ProcessMethod: ProcessDelete,
			Resources: []Resource{
				&testResourcePatchDispatchChange{
					Patch: func() *Patch {
						p := NewPatch()

						p.SetCreateChange(true)

						return p
					}(),
				},
			},
			ErrorMatcher: nil,
		},

		// Case 2
		{
			ProcessMethod: ProcessDelete,
			Resources: []Resource{
				&testResourcePatchDispatchChange{
					Patch: func() *Patch {
						p := NewPatch()

						p.SetDeleteChange(true)

						return p
					}(),
				},
				&testResourcePatchDispatchChange{
					Patch: func() *Patch {
						p := NewPatch()

						p.SetUpdateChange(true)

						return p
					}(),
				},
			},
			ErrorMatcher: nil,
		},

		// Case 3
		{
			ProcessMethod: ProcessDelete,
			Resources: []Resource{
				&testResourcePatchDispatchChange{
					Patch: func() *Patch {
						p := NewPatch()

						p.SetDeleteChange(true)

						return p
					}(),
				},
			},
			ErrorMatcher: nil,
		},

		// Case 4
		{
			ProcessMethod: ProcessDelete,
			Resources: []Resource{
				&testResourcePatchDispatchChange{
					Patch: func() *Patch {
						p := NewPatch()

						p.SetDeleteChange(true)
						p.SetUpdateChange(true)

						return p
					}(),
				},
				&testResourcePatchDispatchChange{
					Patch: func() *Patch {
						p := NewPatch()

						p.SetCreateChange(true)
						p.SetDeleteChange(true)

						return p
					}(),
				},
			},
			ErrorMatcher: nil,
		},

		// Case 5
		{
			ProcessMethod: ProcessDelete,
			Resources: []Resource{
				&testResourcePatchDispatchChange{
					Patch: func() *Patch {
						p := NewPatch()

						p.SetCreateChange(true)
						p.SetDeleteChange(true)
						p.SetUpdateChange(true)

						return p
					}(),
				},
				&testResourcePatchDispatchChange{
					Patch: func() *Patch {
						p := NewPatch()
						return p
					}(),
				},
			},
			ErrorMatcher: nil,
		},
	}

	for i, tc := range testCases {
		err := tc.ProcessMethod(context.TODO(), nil, tc.Resources)
		if err != nil {
			if tc.ErrorMatcher == nil {
				t.Fatal("test", i, "expected", nil, "got", err)
			} else if !tc.ErrorMatcher(err) {
				t.Fatal("test", i, "expected", true, "got", false)
			}
		}
	}
}

type testResourcePatchDispatchChange struct {
	Patch *Patch
}

func (r *testResourcePatchDispatchChange) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	return nil, nil
}

func (r *testResourcePatchDispatchChange) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	return nil, nil
}

func (r *testResourcePatchDispatchChange) NewUpdatePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*Patch, error) {
	return r.Patch, nil
}

func (r *testResourcePatchDispatchChange) NewDeletePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*Patch, error) {
	return r.Patch, nil
}

func (r *testResourcePatchDispatchChange) Name() string {
	return "testResourcePatchDispatchChange"
}

func (r *testResourcePatchDispatchChange) ApplyCreateChange(ctx context.Context, obj, createState interface{}) error {
	createChange, ok := r.Patch.getCreateChange()
	if ok {
		if createChange != createState {
			panic(fmt.Sprintf("expected '%s' got '%s'", createChange, createState))
		}
	}

	return nil
}

func (r *testResourcePatchDispatchChange) ApplyDeleteChange(ctx context.Context, obj, deleteState interface{}) error {
	deleteChange, ok := r.Patch.getDeleteChange()
	if ok {
		if deleteChange != deleteState {
			panic(fmt.Sprintf("expected '%s' got '%s'", deleteChange, deleteState))
		}
	}

	return nil
}

func (r *testResourcePatchDispatchChange) ApplyUpdateChange(ctx context.Context, obj, updateState interface{}) error {
	updateChange, ok := r.Patch.getUpdateChange()
	if ok {
		if updateChange != updateState {
			panic(fmt.Sprintf("expected '%s' got '%s'", updateChange, updateState))
		}
	}

	return nil
}
