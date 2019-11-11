package test

import (
	"context"
)

type NopBasicResource struct {
}

func NewNopBasicResource() *NopBasicResource {
	return &NopBasicResource{}
}

func (*NopBasicResource) Name() string {
	return "NopBasicResource"
}
func (*NopBasicResource) EnsureCreated(context.Context, interface{}) error {
	return nil
}
func (*NopBasicResource) EnsureDeleted(context.Context, interface{}) error {
	return nil
}
