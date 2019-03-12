package kubernetes

import (
	"github.com/giantswarm/operatorkit/flag/service/kubernetes/tls"
	"github.com/giantswarm/operatorkit/flag/service/kubernetes/watch"
)

// Kubernetes is a data structure to hold Kubernetes specific command line
// configuration flags.
type Kubernetes struct {
	Address    string
	InCluster  string
	KubeConfig string
	TLS        tls.TLS
	Watch      watch.Watch
}
