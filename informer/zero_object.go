package informer

import (
	"k8s.io/apimachinery/pkg/runtime"
)

// ZeroObjectFactoryFuncs implements ZeroObjectFactory.
type ZeroObjectFactoryFuncs struct {
	NewObjectFunc     func() runtime.Object
	NewObjectListFunc func() runtime.Object
}

func (z ZeroObjectFactoryFuncs) NewObject() runtime.Object {
	return z.NewObjectFunc()
}

func (z ZeroObjectFactoryFuncs) NewObjectList() runtime.Object {
	return z.NewObjectListFunc()
}
