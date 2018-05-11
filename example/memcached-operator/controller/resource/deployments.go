package resource

import (
	"context"

	examplev1alpha1 "github.com/giantswarm/apiextensions/pkg/apis/example/v1alpha1"
	"github.com/giantswarm/operatorkit/example/memcached-operator/logger"
)

const (
	deploymentsName = "deployments"
)

type DeploymentsConfig struct {
}

type Deployments struct {
}

func NewDeployments(config DeploymentsConfig) (*Deployments, error) {
	d := &Deployments{}
	return d, nil
}
func (d *Deployments) Name() string {
	return deploymentsName
}

func (d *Deployments) EnsureCreated(ctx context.Context, obj interface{}) error {
	_ = obj.(*examplev1alpha1.MemcachedConfig)

	logger.LogCtx(ctx, "level", "debug", "message", "TODO: implement Deployments.EnsureCreated")
	return nil
}

func (d *Deployments) EnsureDeleted(ctx context.Context, obj interface{}) error {
	_ = obj.(*examplev1alpha1.MemcachedConfig)

	logger.LogCtx(ctx, "level", "debug", "message", "TODO: implement Deployments.EnsureDeleted")
	return nil
}
