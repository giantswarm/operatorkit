package nopointer

import "context"

type Handler struct{}

func (h Handler) EnsureCreated(ctx context.Context, obj interface{}) error {
	return nil
}

func (h Handler) EnsureDeleted(ctx context.Context, obj interface{}) error {
	return nil
}
