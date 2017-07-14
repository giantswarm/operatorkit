// Package etcdv2 provides a service that implements an etcd storage.
package etcdv2

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/coreos/etcd/client"

	microerror "github.com/giantswarm/microkit/error"
)

// Config represents the configuration used to create a service.
type Config struct {
	// Dependencies.
	EtcdClient client.Client

	// Settings.
	Prefix string
}

// DefaultConfig provides a default configuration to create a new service by
// best effort.
func DefaultConfig() Config {
	etcdConfig := client.Config{
		Endpoints: []string{"http://127.0.0.1:2379"},
		Transport: client.DefaultTransport,
	}
	etcdClient, err := client.New(etcdConfig)
	if err != nil {
		panic(err)
	}

	return Config{
		// Dependencies.
		EtcdClient: etcdClient,

		// Settings.
		Prefix: "",
	}
}

// New creates a new configured service.
func New(config Config) (*Service, error) {
	// Dependencies.
	if config.EtcdClient == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "etcd client must not be empty")
	}

	newService := &Service{
		// Dependencies.
		etcdClient: config.EtcdClient,

		// Internals.
		keyClient: client.NewKeysAPI(config.EtcdClient),

		// Settings.
		prefix: config.Prefix,
	}

	return newService, nil
}

// Service provides the actual service implementation.
type Service struct {
	// Dependencies.
	etcdClient client.Client

	// Internals.
	keyClient client.KeysAPI

	// Settings.
	prefix string
}

func (s *Service) Create(ctx context.Context, key, value string) error {
	_, err := s.keyClient.Create(ctx, s.key(key), value)
	if IsEtcdKeyAlreadyExists(err) {
		return nil
	} else if err != nil {
		return microerror.MaskAny(err)
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, key string) error {
	options := &client.DeleteOptions{
		Recursive: true,
	}
	_, err := s.keyClient.Delete(ctx, s.key(key), options)
	if client.IsKeyNotFound(err) {
		return microerror.MaskAnyf(notFoundError, "%s", err)
	} else if err != nil {
		return microerror.MaskAny(err)
	}

	return nil
}

func (s *Service) Exists(ctx context.Context, key string) (bool, error) {
	options := &client.GetOptions{
		Quorum: true,
	}
	_, err := s.keyClient.Get(ctx, s.key(key), options)
	if client.IsKeyNotFound(err) {
		return false, nil
	} else if err != nil {
		return false, microerror.MaskAny(err)
	}

	return true, nil
}

func (s *Service) List(ctx context.Context, key string) ([]string, error) {
	options := &client.GetOptions{
		Recursive: true,
		Quorum:    true,
	}
	resp, err := s.keyClient.Get(ctx, s.key(key), options)
	if client.IsKeyNotFound(err) {
		return nil, microerror.MaskAny(notFoundError)
	} else if err != nil {
		return nil, microerror.MaskAny(err)
	}
	if resp.Node == nil || resp.Node.Dir == false {
		return nil, microerror.MaskAny(notFoundError)
	}

	var children []string

	for _, node := range resp.Node.Nodes {
		if node.Dir == true {
			continue
		}
		if !strings.HasPrefix(node.Key, s.key(key)) {
			return nil, microerror.MaskAny(notFoundError)
		}
		children = append(children, node.Key[len(s.key(key))+1:])
	}

	if len(children) == 0 {
		return nil, microerror.MaskAny(notFoundError)
	}

	return children, nil
}

func (s *Service) Search(ctx context.Context, key string) (string, error) {
	options := &client.GetOptions{
		Quorum: true,
	}
	clientResponse, err := s.keyClient.Get(ctx, s.key(key), options)
	if client.IsKeyNotFound(err) {
		return "", microerror.MaskAnyf(notFoundError, key)
	} else if err != nil {
		return "", microerror.MaskAny(err)
	}

	return clientResponse.Node.Value, nil
}

func (s *Service) key(key string) string {
	return filepath.Clean(filepath.Join("/", s.prefix, key))
}
