package etcd

import (
	"testing"
)

func Test_Prefix_key(t *testing.T) {
	config := DefaultConfig()
	config.Prefix = "foo"
	newStorage, err := New(config)
	if err != nil {
		panic(err)
	}

	key := newStorage.key("bar")
	if key != "/foo/bar" {
		t.Fatal("expected", "/foo/bar", "got", key)
	}
}

func Test_Prefix_key_slash(t *testing.T) {
	config := DefaultConfig()
	config.Prefix = "/foo/"
	newStorage, err := New(config)
	if err != nil {
		panic(err)
	}

	key := newStorage.key("bar")
	if key != "/foo/bar" {
		t.Fatal("expected", "/foo/bar", "got", key)
	}
}
