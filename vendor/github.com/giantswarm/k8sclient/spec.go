package k8sclient

import (
	"github.com/giantswarm/apiextensions/pkg/clientset/versioned"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/giantswarm/k8sclient/k8scrdclient"
)

type Interface interface {
	CRDClient() k8scrdclient.Interface
	CtrlClient() client.Client
	DynClient() dynamic.Interface
	ExtClient() apiextensionsclient.Interface
	G8sClient() versioned.Interface
	K8sClient() kubernetes.Interface
	RESTClient() rest.Interface
	RESTConfig() *rest.Config
	Scheme() *runtime.Scheme
}
