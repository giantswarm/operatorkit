package framework

import "context"

func DefaultResourceRouter(rs []Resource) func(ctx context.Context, obj interface{}) ([]Resource, error) {
	return func(ctx context.Context, obj interface{}) ([]Resource, error) {
		return rs, nil
	}
}
