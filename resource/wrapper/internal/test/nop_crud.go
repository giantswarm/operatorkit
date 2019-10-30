package test

import (
	"context"

	"github.com/giantswarm/operatorkit/resource/crud"
)

type NopCRUD struct {
}

func NewNopCRUD() *NopCRUD {
	return &NopCRUD{}
}

func (*NopCRUD) Name() string {
	return "NopCRUD"
}
func (*NopCRUD) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	return nil, nil
}
func (*NopCRUD) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	return nil, nil
}
func (*NopCRUD) NewUpdatePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*crud.Patch, error) {
	return nil, nil
}
func (*NopCRUD) NewDeletePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*crud.Patch, error) {
	return nil, nil
}
func (*NopCRUD) ApplyCreateChange(ctx context.Context, obj, createChange interface{}) error {
	return nil
}
func (*NopCRUD) ApplyDeleteChange(ctx context.Context, obj, deleteChange interface{}) error {
	return nil
}
func (*NopCRUD) ApplyUpdateChange(ctx context.Context, obj, updateChange interface{}) error {
	return nil
}
