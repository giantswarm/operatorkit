// +build k8srequired

package client

import "github.com/giantswarm/operatorkit/framework"

type Interface interface {
	// CRUD functions for objects.
	CreateObject(namespace string, obj interface{}) (interface{}, error)
	DeleteObject(name, namespace string) error
	GetObject(name, namespace string) (interface{}, error)
	UpdateObject(namespace string, obj interface{}) (interface{}, error)

	// Functions for test setup and teardown.
	MustSetup(namespace string)
	MustTeardown(namespace string)

	// Getters and setters.
	GetFramework() *framework.Framework
}

type Config struct {
	Resources []framework.Resource

	Name      string
	Namespace string
}
