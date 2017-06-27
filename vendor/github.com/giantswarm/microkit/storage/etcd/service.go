// Package etcd provides a service that implements an etcd storage, using the etcd v3 protocol.
package etcd

import (
	"path/filepath"

	"github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"

	microerror "github.com/giantswarm/microkit/error"
)

// Config represents the configuration used to create a etcd service.
type Config struct {
	// Dependencies.
	EtcdClient *clientv3.Client

	// Settings.
	Prefix string
}

// DefaultConfig provides a default configuration to create a new etcd service
// by best effort.
func DefaultConfig() Config {
	etcdConfig := clientv3.Config{
		Endpoints: []string{"http://127.0.0.1:2379"},
	}
	etcdClient, err := clientv3.New(etcdConfig)
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

// New creates a new configured etcd service.
func New(config Config) (*Service, error) {
	// Dependencies.
	if config.EtcdClient == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "etcd client must not be empty")
	}

	newService := &Service{
		// Dependencies.
		etcdClient: config.EtcdClient,

		// Internals.
		keyClient: clientv3.NewKV(config.EtcdClient),

		// Settings.
		prefix: config.Prefix,
	}

	return newService, nil
}

// Service is the etcd service.
type Service struct {
	// Dependencies.
	etcdClient *clientv3.Client

	// Internals.
	keyClient clientv3.KV

	// Settings.
	prefix string
}

func (s *Service) Create(ctx context.Context, key, value string) error {
	key = s.key(key)

	_, err := s.keyClient.Put(ctx, key, value)
	if err != nil {
		return microerror.MaskAny(err)
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, key string) error {
	key = s.key(key)

	resp, err := s.keyClient.Delete(ctx, key)
	if err != nil {
		return microerror.MaskAny(err)
	}
	if resp.Deleted == 0 {
		return microerror.MaskAny(notFoundError)
	}

	return nil
}

func (s *Service) Exists(ctx context.Context, key string) (bool, error) {
	_, err := s.Search(ctx, key)
	if IsNotFound(err) {
		return false, nil
	} else if err != nil {
		return false, microerror.MaskAny(err)
	}

	return true, nil
}

func (s *Service) List(ctx context.Context, key string) ([]string, error) {
	opts := []clientv3.OpOption{
		clientv3.WithKeysOnly(),
		clientv3.WithPrefix(),
	}

	key = s.key(key)
	res, err := s.keyClient.Get(ctx, key, opts...)
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	if res.Count == 0 {
		return nil, microerror.MaskAnyf(notFoundError, key)
	}

	var list []string

	i := len(key)
	for _, kv := range res.Kvs {
		k := string(kv.Key)

		if len(k) <= i+1 {
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
	res, err := s.keyClient.Get(ctx, s.key(key))
	if err != nil {
		return "", microerror.MaskAny(err)
	}

	if res.Count == 0 {
		return "", microerror.MaskAnyf(notFoundError, key)
	}

	if res.Count > 1 {
		return "", microerror.MaskAnyf(multipleValuesError, key)
	}

	return string(res.Kvs[0].Value), nil
}

func (s *Service) key(key string) string {
	return filepath.Clean(filepath.Join("/", s.prefix, key))
}
