// Package memory provides a service that implements a memory storage.
package memory

import (
	"strings"
	"sync"

	"golang.org/x/net/context"

	microerror "github.com/giantswarm/microkit/error"
)

// Config represents the configuration used to create a memory service.
type Config struct {
}

// DefaultConfig provides a default configuration to create a new memory
// service by best effort.
func DefaultConfig() Config {
	return Config{}
}

// New creates a new configured memory service.
func New(config Config) (*Service, error) {
	newService := &Service{
		storage: map[string]string{},
		mutex:   sync.Mutex{},
	}

	return newService, nil
}

// Service is the memory service.
type Service struct {
	// Internals.
	storage map[string]string
	mutex   sync.Mutex
}

func (s *Service) Create(ctx context.Context, key, value string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.storage[key] = value

	return nil
}

func (s *Service) Delete(ctx context.Context, key string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.storage, key)

	return nil
}

func (s *Service) Exists(ctx context.Context, key string) (bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, ok := s.storage[key]

	return ok, nil
}

func (s *Service) List(ctx context.Context, key string) ([]string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var list []string

	i := len(key)
	for k, _ := range s.storage {
		if len(k) <= i+1 {
			continue
		}
		if !strings.HasPrefix(k, key) {
			continue
		}

		if k[i] != '/' {
			// We want to ignore all keys that are not separated by slash. When there
			// is a key stored like "foo/bar/baz", listing keys using "foo/ba" should
			// not succeed.
			continue
		}

		list = append(list, k[i+1:])
	}

	if len(list) == 0 {
		return nil, microerror.MaskAnyf(notFoundError, key)
	}

	return list, nil
}

func (s *Service) Search(ctx context.Context, key string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	value, ok := s.storage[key]
	if ok {
		return value, nil
	}

	return "", microerror.MaskAnyf(notFoundError, key)
}
