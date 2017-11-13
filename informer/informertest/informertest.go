package informertest

import (
	"context"

	"k8s.io/apimachinery/pkg/watch"
)

type InformerTest struct{}

func New() *InformerTest {
	return &InformerTest{}
}

func (i *InformerTest) Watch(ctx context.Context) (chan watch.Event, chan watch.Event, chan error) {
	return nil, nil, nil
}
