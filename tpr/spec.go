package tpr

import (
	"k8s.io/apimachinery/pkg/runtime"
)

type Interface interface {
	CreateAndWait() error
}

// ZeroObjectFuncs provides zero values of an object and objects' list ready to
// be decoded. The provided zero values must not be reused by zeroObjectFactory.
type ZeroObjectFactory interface {
	NewObject() runtime.Object
	NewObjectList() runtime.Object
}
