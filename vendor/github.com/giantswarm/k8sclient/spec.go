package k8sclient

import (
	"github.com/giantswarm/apiextensions/pkg/clientset/versioned"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Interface interface {
	DynClient() dynamic.Interface
	ExtClient() apiextensionsclient.Interface
	G8sClient() versioned.Interface
	K8sClient() kubernetes.Interface
	RESTClient() rest.Interface
	RESTConfig() *rest.Config
}
