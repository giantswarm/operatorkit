package wrapper

import "github.com/giantswarm/operatorkit/controller"

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
	Controller() *controller.Controller
}
