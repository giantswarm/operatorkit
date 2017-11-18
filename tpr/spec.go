package tpr

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
)

type Interface interface {
	CreateAndWait() error
	CollectMetrics(context.Context)
}

// ZeroObjectFuncs provides zero values of an object and objects' list ready to
// be decoded. The provided zero values must not be reused by zeroObjectFactory.
type ZeroObjectFactory interface {
	NewObject() runtime.Object
	NewObjectList() runtime.Object
}

type TPOList struct {
	Items []interface{} `json:"items"`
}
