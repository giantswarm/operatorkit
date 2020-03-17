package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// CommonClusterObject represents common interface for all provider specific
// cluster objects.
type CommonClusterObject interface {
	metav1.Object
	runtime.Object
	CommonClusterStatusGetSetter
}

// CommonClusterStatusGetSetter provides abstract way to manipulate common
// provider independent cluster status field in provider CR's status.
type CommonClusterStatusGetSetter interface {
	GetCommonClusterStatus() CommonClusterStatus
	SetCommonClusterStatus(ccs CommonClusterStatus)
}
