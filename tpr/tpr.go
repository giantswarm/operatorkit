package tpr

import (
	"time"

	"github.com/cenkalti/backoff"
	microerror "github.com/giantswarm/microkit/error"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/errors"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

const (
	// ResyncPeriod is the interval at which the Informer cache is invalidated,
	// and the lister function is called.
	ResyncPeriod = 1 * time.Minute

	tprInitMaxElapsedTime = 2 * time.Minute
)

// Config is a TPR configuration.
type Config struct {
	// Dependencies.
	Clientset kubernetes.Interface

	// Settings.

	// Name takes the form <kind>.<group> (Note: The group is also called
	// a domain). You are expected to provide a unique kind and group name
	// in order to avoid conflicts with other ThirdPartyResource objects.
	// Kind names will be converted to CamelCase when creating instances of
	// the ThirdPartyResource.  Hyphens in the kind are assumed to be word
	// breaks. For instance the kind camel-case would be converted to
	// CamelCase but camelcase would be converted to Camelcase.
	Name string

	// Version is TPR version, e.g. v1.
	Version string

	// Description is free text description.
	Description string
}

// TPR allows easy operations on ThirdPartyResources. See
// https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-third-party-resource/
// for details.
type TPR struct {
	clientset kubernetes.Interface

	name        string // see Config.Name
	kind        string // see Config.Name
	group       string // see Config.Name
	version     string
	apiVersion  string // apiVersion is group/version
	description string

	// API for this TPR kind name.
	resourceName string
}

// New creates a new TPR.
func New(config Config) (*TPR, error) {
	if config.Clientset == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "k8s clientset must be set")
	}
	if config.Name == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "name must not be empty")
	}
	if config.Version == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "version must not be empty")
	}
	if config.Description == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "description must not be empty")
	}

	kind, group, err := extractKindAndGroup(config.Name)
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	tpr := &TPR{
		clientset: config.Clientset,

		name:        config.Name,
		kind:        kind,
		group:       group,
		version:     config.Version,
		apiVersion:  group + "/" + config.Version,
		description: config.Description,

		resourceName: unsafeGuessKindToResource(kind),
	}
	return tpr, nil
}

// Kind returns a TPR kind extracted from Name. Useful when creating
// ThirdPartyObjects. See Config.Name godoc for details.
func (t *TPR) Kind() string { return t.kind }

// APIVersion returns a TPR APIVersion created from group and version. It takes
// format <group>/<version>. Useful for creating ThirdPartyObjects. See
// Config.Name and Config.Version for details.
func (t *TPR) APIVersion() string { return t.apiVersion }

// Name returns a TPR name provided with Config.Name.
func (t *TPR) Name() string { return t.name }

// Group returns a TPR group extracted from Name. See Config.Name godoc for details.
func (t *TPR) Group() string { return t.group }

// Endpoint returns a TPR resource endpoint registered in the Kubernetes API
// under a given namespace. The default namespace will be used when the
// argument is an empty string.
func (t *TPR) Endpoint(namespace string) string {
	nsResource := t.resourceName
	if len(namespace) != 0 {
		nsResource = "namespace/" + namespace + "/" + t.resourceName
	}
	return "/apis/" + t.group + "/" + t.version + "/" + nsResource
}

// WatchEndpoint returns a TPR watch resource endpoint registered in the
// Kubernetes API under a given namespace. The default namespace will be used
// when the argument is an empty string.
func (t *TPR) WatchEndpoint(namespace string) string {
	nsResource := "watch/" + t.resourceName
	if len(namespace) != 0 {
		nsResource = "namespace/" + namespace + "/watch/" + t.resourceName
	}
	return "/apis/" + t.group + "/" + t.version + "/" + nsResource
}

// CreateAndWait creates a TPR and waits till it is initialized in the cluster.
// Retruns alreadyExistsError when the resource already exists.
func (t *TPR) CreateAndWait() error {
	initBackOff := backoff.NewExponentialBackOff()
	initBackOff.MaxElapsedTime = tprInitMaxElapsedTime
	return t.CreateAndWaitBackOff(initBackOff)
}

// CreateAndWaitBackOff creates a TPR and waits till it is initialized in the
// cluster. It allows to pass a custom initialization back off policy used to
// poll for TPR readiness. Retruns alreadyExistsError when the resource already
// exists.
func (t *TPR) CreateAndWaitBackOff(initBackOff backoff.BackOff) error {
	err := t.create()
	if err != nil {
		return microerror.MaskAnyf(err, "creating TPR %s", t.name)
	}
	err = t.waitInit(initBackOff)
	if err != nil {
		return microerror.MaskAnyf(err, "waiting for TPR %s initialization", t.name)
	}
	return nil
}

// create is extracted for testing because fake REST client does not work.
// Therefore waitInit can not be tested.
func (t *TPR) create() error {
	tpr := &v1beta1.ThirdPartyResource{
		ObjectMeta: v1.ObjectMeta{
			Name: t.name,
		},
		Versions: []v1beta1.APIVersion{
			{Name: t.version},
		},
		Description: t.description,
	}

	_, err := t.clientset.ExtensionsV1beta1().ThirdPartyResources().Create(tpr)
	if err != nil && errors.IsAlreadyExists(err) {
		return microerror.MaskAny(alreadyExistsError)
	}
	if err != nil {
		return microerror.MaskAnyf(err, "creating TPR %s", t.name)
	}
	return nil
}

func (t *TPR) waitInit(retry backoff.BackOff) error {
	endpoint := t.Endpoint("")
	op := func() error {
		_, err := t.clientset.CoreV1().RESTClient().Get().RequestURI(endpoint).DoRaw()
		return err
	}

	err := backoff.Retry(op, retry)

	if errors.IsNotFound(err) {
		err = tprInitTimeoutError
	}
	return microerror.MaskAnyf(err, "requesting TPR %s", t.name)
}
