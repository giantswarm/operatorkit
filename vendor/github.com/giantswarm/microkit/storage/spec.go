package storage

import (
	"context"
)

// Service represents the abstraction for underlying storage backends. A storage
// service implementation does not care about specific data types. All the
// storage cares about are key-value pairs. Services making use of the storage
// have to take care about specific types they care about them self.
type Service interface {
	// Create stores the given value under the given key. Keys and values might
	// have specific schemes depending on the specific storage implementation.
	// E.g. an etcd storage implementation will allow keys to be defined as paths:
	// path/to/key. Values might be JSON strings in case the service using the
	// storage wants to store its data as JSON strings.
	Create(ctx context.Context, key, value string) error
	// Delete removes the value stored under the given key.
	Delete(ctx context.Context, key string) error
	// Exists checks if a value under the given key exists or not.
	Exists(ctx context.Context, key string) (bool, error)
	// List does a lookup for all values stored under key and returns them, if
	// any.
	List(ctx context.Context, key string) ([]string, error)
	// Search does a lookup for the value stored under key and returns it, if any.
	Search(ctx context.Context, key string) (string, error)
}
