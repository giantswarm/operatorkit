package v1alpha1

import (
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

const (
	crDocsAnnotation         = "giantswarm.io/docs"
	kindRelease              = "Release"
	releaseDocumentationLink = "https://pkg.go.dev/github.com/giantswarm/apiextensions/pkg/apis/release/v1alpha1?tab=doc#Release"
	releaseCRDYAML           = `apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: releases.release.giantswarm.io
spec:
  additionalPrinterColumns:
    - name: Kubernetes version
      type: string
      description: Version of the kubernetes component in this release
      JSONPath: .spec.components[?(@.name=="kubernetes")].version
    - name: State
      type: string
      description: State of the release
      JSONPath: .spec.state
    - name: Age
      type: date
      description: Time since release creation
      JSONPath: .spec.date
  group: release.giantswarm.io
  names:
    kind: Release
    plural: releases
    shortNames:
    - rel
    singular: release
  preserveUnknownFields: false
  scope: Cluster
  validation:
    openAPIV3Schema:
      description: |
        A Release holds information about a particular version of the Giant Swarm platform which
        can be used as a target for creation or upgrade of a cluster. It is a tested package
        comprising a particular Kubernetes version along with compatible Giant Swarm operators,
        monitoring, and default apps.
      properties:
        metadata:
          properties:
            name:
              pattern: ^v(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)$
              type: string
          type: object
        spec:
          description: |
            Spec holds the data defining the desired state of a release.
          properties:
            date:
              type: string
              format: date-time
            apps:
              description: |
                Apps is a list of Giant Swarm-managed apps which will be installed by default
                on clusters created with this release version.
              items:
                properties:
                  componentVersion:
                    description: |
                      Component version is the upstream version of this app. It may be empty if this
                      is a Giant Swarm developed app.
                    pattern: ^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$
                    type: string
                  name:
                    description: |
                      Name is the name of the app.
                    minLength: 1
                    type: string
                  version:
                    description: |
                      Version is the internal version of the app managed by Giant Swarm. Because apps
                      may be released without upstream changes, this will generally differ from the
                      component version.
                    pattern: ^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$
                    type: string
                required:
                - name
                - version
                type: object
              type: array
            components:
              description: |
                Components is a list of internal and upstream components making up the core of the cluster.
              items:
                properties:
                  name:
                    description: |
                      Name is the name of the component.
                    minLength: 1
                    type: string
                  version:
                    description: |
                      Version is the semantic version of the component.
                    pattern: ^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$
                    type: string
                required:
                - name
                - version
                type: object
              minItems: 1
              type: array
            state:
              description: |
                State indicates how this release should be used. "wip" means the release is a work
                in progress. "deprecated" means old clusters using this version will continue to function
                but new clusters should use a more recent release. "active" means this is a current
                supported release.
              pattern: ^(active|deprecated|wip)$
              type: string
          required:
          - components
          - apps
          - state
          - date
          type: object
      required:
      - metadata
      type: object
  versions:
  - name: v1alpha1
    served: true
    storage: true
`
)

type ReleaseState string

var (
	stateActive     ReleaseState = "active"
	stateDeprecated ReleaseState = "deprecated"
	stateWIP        ReleaseState = "wip"
	releaseCRD      *apiextensionsv1beta1.CustomResourceDefinition
)

func init() {
	err := yaml.UnmarshalStrict([]byte(releaseCRDYAML), &releaseCRD)
	if err != nil {
		panic(err)
	}
}

// NewReleaseCRD returns a new custom resource definition for Release.
func NewReleaseCRD() *apiextensionsv1beta1.CustomResourceDefinition {
	return releaseCRD.DeepCopy()
}

func NewReleaseTypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{
		APIVersion: SchemeGroupVersion.String(),
		Kind:       kindRelease,
	}
}

func NewReleaseCR() *Release {
	return &Release{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				crDocsAnnotation: releaseDocumentationLink,
			},
		},
		TypeMeta: NewReleaseTypeMeta(),
	}
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Release is a Kubernetes resource (CR) which is based on the Release CRD defined above.
//
// An example Release resource can be viewed here
// https://github.com/giantswarm/apiextensions/blob/master/docs/cr/release.giantswarm.io_v1alpha1_release.yaml
type Release struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec              ReleaseSpec `json:"spec" yaml:"spec"`
}

type ReleaseSpec struct {
	// Apps describes apps used in this release.
	Apps []ReleaseSpecApp `json:"apps" yaml:"apps"`
	// Components describes components used in this release.
	Components []ReleaseSpecComponent `json:"components" yaml:"components"`
	// Date that the release became active.
	Date *DeepCopyTime `json:"date" yaml:"date"`
	// State indicates the availability of the release: deprecated, active, or wip.
	State ReleaseState `json:"state" yaml:"state"`
}

type ReleaseSpecComponent struct {
	// Name of the component.
	Name string `json:"name" yaml:"name"`
	// Version of the component.
	Version string `json:"version" yaml:"version"`
}

type ReleaseSpecApp struct {
	// Version of the upstream component used in the app.
	ComponentVersion string `json:"componentVersion,omitempty" yaml:"componentVersion,omitempty"`
	// Name of the app.
	Name string `json:"name" yaml:"name"`
	// Version of the app.
	Version string `json:"version" yaml:"version"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ReleaseList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata" yaml:"metadata"`
	Items           []Release `json:"items" yaml:"items"`
}
