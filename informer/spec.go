package informer

import (
	"k8s.io/apimachinery/pkg/runtime"
)

// WatchEndpointProvider knows how to construct watch endpoints. E.g. CRD
// implements WatchEndpointProvider.
type WatchEndpointProvider interface {
	WatchEndpoint() string
}

// ZeroObjectFuncs provides zero values of an object and objects' list ready to
// be decoded. The provided zero values must not be reused by zeroObjectFactory.
type ZeroObjectFactory interface {
	NewObject() runtime.Object
	NewObjectList() runtime.Object
}
