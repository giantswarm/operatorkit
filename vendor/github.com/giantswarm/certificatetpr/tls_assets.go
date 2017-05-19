package certificatetpr

// ClusterComponent represents the individual component of a k8s cluster, e.g.
// the API server, or etcd These are used when getting a secret from the k8s
// API, to identify the component the secret belongs to.
type ClusterComponent string

func (c ClusterComponent) String() string {
	return string(c)
}

// TLSAssetType represents the type of TLS asset, e.g. a CA certificate, or a
// certificate key These are used when getting a secret from the k8s API, to
// identify the specific type of TLS asset that is contained in the secret.
type TLSAssetType string

func (t TLSAssetType) String() string {
	return string(t)
}

// These constants are used to match each asset in the secret.
const (
	// CA is the key for the CA certificate.
	CA TLSAssetType = "ca"
	// Crt is the key for the certificate.
	Crt TLSAssetType = "crt"
	// Key is the key for the key.
	Key TLSAssetType = "key"
)

// These constants are used to match different components of the cluster when
// parsing a secret received from the API.
const (
	// APIComponent is the API server component.
	APIComponent ClusterComponent = "api"
	// WorkerComponent is a worker component.
	WorkerComponent ClusterComponent = "worker"
	// EtcdComponent is the etcd cluster component.
	EtcdComponent ClusterComponent = "etcd"
	// CalicoComponent is the calico component.
	CalicoComponent ClusterComponent = "calico"
	// ServiceAccountComponent is the service-account component.
	ServiceAccountComponent ClusterComponent = "service-account"
)

// These constants are used when filtering the secrets, to only retrieve the
// ones we are interested in.
const (
	// ComponentLabel is the label used in the secret to identify a cluster
	// component.
	ComponentLabel string = "clusterComponent"
	// ClusterIDLabel is the label used in the secret to identify a cluster.
	ClusterIDLabel string = "clusterID"
)

// AssetsBundleKey is a struct key for an AssetsBundle cfr.
// https://blog.golang.org/go-maps-in-action
type AssetsBundleKey struct {
	Component ClusterComponent
	Type      TLSAssetType
}

// AssetsBundle is a structure that contains all the assets for all the
// components.
type AssetsBundle map[AssetsBundleKey][]byte

// ClusterComponents is a slice enumerating all the components that make up the
// cluster.
var ClusterComponents = []ClusterComponent{
	APIComponent,
	WorkerComponent,
	EtcdComponent,
	CalicoComponent,
	ServiceAccountComponent,
}

// TLSAssetTypes is a slice enumerating all the TLS assets we need to boot the
// cluster.
var TLSAssetTypes = []TLSAssetType{CA, Crt, Key}

// ValidComponent looks for el among the components.
func ValidComponent(el ClusterComponent, components []ClusterComponent) bool {
	for _, v := range components {
		if el == v {
			return true
		}
	}
	return false
}

// CompactTLSAssets is a struct used by operators to store stringified TLS assets.
type CompactTLSAssets struct {
	APIServerCA       string
	APIServerKey      string
	APIServerCrt      string
	WorkerCA          string
	WorkerKey         string
	WorkerCrt         string
	ServiceAccountCA  string
	ServiceAccountKey string
	ServiceAccountCrt string
	CalicoClientCA    string
	CalicoClientKey   string
	CalicoClientCrt   string
	EtcdServerCA      string
	EtcdServerKey     string
	EtcdServerCrt     string
}
