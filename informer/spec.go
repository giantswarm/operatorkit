package informer

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
)

// WatcherFactory is able to create watchers on demand. It takes a watch
// endpoint and a ZeroObjectFactory to be able to decode watched events.
type WatcherFactory func() (watch.Interface, error)

// ZeroObjectFuncs provides zero values of an object and objects' list ready to
// be decoded. The provided zero values must not be reused by zeroObjectFactory.
type ZeroObjectFactory interface {
	NewObject() runtime.Object
	NewObjectList() runtime.Object
}
