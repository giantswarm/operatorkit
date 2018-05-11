package resource

import (
	"context"

	examplev1alpha1 "github.com/giantswarm/apiextensions/pkg/apis/example/v1alpha1"
	"github.com/giantswarm/operatorkit/example/memcached-operator/logger"
)

const (
	servicesName = "services"
)

type ServicesConfig struct {
}

type Services struct {
}

func NewServices(config ServicesConfig) (*Services, error) {
	s := &Services{}
	return s, nil
}
func (s *Services) Name() string {
	return servicesName
}

func (s *Services) EnsureCreated(ctx context.Context, obj interface{}) error {
	_ = obj.(*examplev1alpha1.MemcachedConfig)

	logger.LogCtx(ctx, "level", "debug", "message", "TODO: implement Services.EnsureCreated")
	return nil
}

func (s *Services) EnsureDeleted(ctx context.Context, obj interface{}) error {
	_ = obj.(*examplev1alpha1.MemcachedConfig)

	logger.LogCtx(ctx, "level", "debug", "message", "TODO: implement Services.EnsureDeleted")
	return nil
}
