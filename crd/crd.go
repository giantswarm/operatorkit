package crd

import (
	"path/filepath"

	"github.com/giantswarm/microerror"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// K8sAPIBase is the base of the Kubernetes API root when defining an absolute
	// path for a rest client request.
	K8sAPIBase = "apis"
)

// Config represents the configuration used to create a new CRD.
type Config struct {
	// Settings.
	Group    string
	Kind     string
	Name     string
	Plural   string
	Singular string
	Scope    string
	Version  string
}

// DefaultConfig provides a default configuration to create a new CRD by best
// effort.
func DefaultConfig() Config {
	return Config{
		// Settings.
		Group:    "",
		Kind:     "",
		Name:     "",
		Plural:   "",
		Singular: "",
		Scope:    "",
		Version:  "",
	}
}

type CRD struct {
	// Settings.
	group    string
	kind     string
	name     string
	plural   string
	singular string
	scope    string
	version  string
}

// New creates a new CRD.
func New(config Config) (*CRD, error) {
	// Settings.
	if config.Group == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Group must not be empty")
	}
	if config.Kind == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Kind must not be empty")
	}
	if config.Name == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Name must not be empty")
	}
	if config.Plural == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Plural must not be empty")
	}
	if config.Singular == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Singular must not be empty")
	}
	if config.Scope == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Scope must not be empty")
	}
	if config.Scope != string(apiextensionsv1beta1.ClusterScoped) && config.Scope != string(apiextensionsv1beta1.NamespaceScoped) {
		return nil, microerror.Maskf(invalidConfigError, "config.Scope must either be '%s' or '%s'", apiextensionsv1beta1.ClusterScoped, apiextensionsv1beta1.NamespaceScoped)
	}
	if config.Version == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Version must not be empty")
	}

	newCRD := &CRD{
		// Settings.
		group:    config.Group,
		kind:     config.Kind,
		name:     config.Name,
		plural:   config.Plural,
		singular: config.Singular,
		scope:    config.Scope,
		version:  config.Version,
	}

	return newCRD, nil
}

func (c *CRD) CreateEndpoint() string {
	return filepath.Join(K8sAPIBase, c.Group(), c.Version(), c.Plural())
}

func (c *CRD) Group() string {
	return c.group
}

func (c *CRD) Kind() string {
	return c.kind
}

func (c *CRD) ListEndpoint() string {
	return filepath.Join(K8sAPIBase, c.Group(), c.Version(), c.Plural())
}

func (c *CRD) Name() string {
	return c.name
}

// NewResource returns a new resource object of the CRD. This might look
// something like the following.
//
//     apiVersion: apiextensions.k8s.io/v1beta1
//     kind: CustomResourceDefinition
//     metadata:
//       name: tests.example.com
//     spec:
//       group: example.com
//       version: v1
//       scope: Cluster
//       names:
//         plural: tests
//         singular: test
//         kind: Test
//
func (c *CRD) NewResource() *apiextensionsv1beta1.CustomResourceDefinition {
	return &apiextensionsv1beta1.CustomResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			Kind:       "CustomResourceDefinition",
			APIVersion: apiextensionsv1beta1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: c.Name(),
		},
		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
			Group:   c.Group(),
			Version: c.Version(),
			Scope:   c.Scope(),
			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
				Kind:     c.Kind(),
				Plural:   c.Plural(),
				Singular: c.Singular(),
			},
		},
	}
}

func (c *CRD) Plural() string {
	return c.plural
}

func (c *CRD) ResourceEndpoint(name string) string {
	return filepath.Join(K8sAPIBase, c.Group(), c.Version(), c.Plural(), name)
}

func (c *CRD) Singular() string {
	return c.singular
}

func (c *CRD) Scope() apiextensionsv1beta1.ResourceScope {
	return apiextensionsv1beta1.ResourceScope(c.scope)
}

func (c *CRD) Version() string {
	return c.version
}

func (c *CRD) WatchEndpoint() string {
	return filepath.Join(K8sAPIBase, c.Group(), c.Version(), "watch", c.Plural())
}
