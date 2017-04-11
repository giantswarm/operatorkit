// +build integration

package etcd

import (
	"testing"

	"golang.org/x/net/context"
)

func Test_CreateExistsSearch(t *testing.T) {
	config := DefaultConfig()
	config.Prefix = "foo"
	newStorage, err := New(config)
	if err != nil {
		panic(err)
	}

	key := "my-key"
	val := "my-val"

	// There should be no key/value pair being stored initially.
	ok, err := newStorage.Exists(context.TODO(), key)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
	if ok {
		t.Fatal("expected", false, "got", true)
	}

	// Creating the key/value pair should work.
	err = newStorage.Create(context.TODO(), key, val)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	// There should be the created key/value pair.
	ok, err = newStorage.Exists(context.TODO(), key)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
	if !ok {
		t.Fatal("expected", true, "got", false)
	}
}

func Test_List(t *testing.T) {
	config := DefaultConfig()
	config.Prefix = "foo"
	newStorage, err := New(config)
	if err != nil {
		panic(err)
	}

	val := "my-val"

	err = newStorage.Create(context.TODO(), "key/one", val)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
	err = newStorage.Create(context.TODO(), "key/two", val)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	values, err := newStorage.List(context.TODO(), "key")
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
	if len(values) != 2 {
		t.Fatal("expected", 2, "got", len(values))
	}
	if values[0] != "one" {
		t.Fatal("expected", "one", "got", values[0])
	}
	if values[1] != "two" {
		t.Fatal("expected", "two", "got", values[1])
	}
}

func Test_List_Nesting(t *testing.T) {
	config := DefaultConfig()
	config.Prefix = "foo"
	newStorage, err := New(config)
	if err != nil {
		panic(err)
	}

	val := "my-val"

	err = newStorage.Create(context.TODO(), "tokend/token/some/scope/id1", val)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
	err = newStorage.Create(context.TODO(), "tokend/token/some/other/scope/id34", val)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	values, err := newStorage.List(context.TODO(), "tokend/token")
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
	if len(values) != 2 {
		t.Fatal("expected", 2, "got", len(values))
	}
	if values[0] != "some/other/scope/id34" {
		t.Fatal("expected", "some/other/scope/id34", "got", values[0])
	}
	if values[1] != "some/scope/id1" {
		t.Fatal("expected", "some/scope/id1", "got", values[1])
	}
}

func Test_List_Invalid_Key(t *testing.T) {
	config := DefaultConfig()
	config.Prefix = "foo"
	newStorage, err := New(config)
	if err != nil {
		panic(err)
	}

	val := "my-val"

	err = newStorage.Create(context.TODO(), "key/one", val)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
	err = newStorage.Create(context.TODO(), "key/two", val)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	_, err = newStorage.List(context.TODO(), "ke")
	if !IsNotFound(err) {
		t.Fatal("expected", true, "got", false)
	}
}
