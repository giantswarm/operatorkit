package informer

import (
	"k8s.io/client-go/pkg/runtime"
)

// Observer functions are called when the cache.ListWatch List and Watch
// functions are called.
type Observer interface {
	OnList()
	OnWatch()
}

// ZeroObjectFuncs provides zero values of an object and objects' list ready to
// be decoded. The provided zero values must not be reused by zeroObjectFactory.
type ZeroObjectFactory interface {
	NewObject() runtime.Object
	NewObjectList() runtime.Object
}
