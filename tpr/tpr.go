package tpr

import (
	"fmt"
	"time"

	microerror "github.com/giantswarm/microkit/error"
	"github.com/giantswarm/operatorkit/util"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/errors"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/rest"
)

const (
	tprInitRetries    = 30
	tprInitRetryDelay = 3 * time.Second
)

type Config struct {
	// Dependencies.
	Clientset kubernetes.Interface
	Rest      rest.Interface

	// Settings.

	// Name of the kind of ThirdPartyObjects. It should be in lower case
	// and hyphen delimited. Kind names will be converted to CamelCase
	// when creating ThirdPartyObjects. Hyphens in the kind are assumed to
	// be word breaks. For instance the kind camel-case would be converted
	// to CamelCase but camelcase would be converted to Camelcase.
	Name string

	// Domain of ThirdPartyResource, e.g. example.com. Along with Name must
	// create an unique pair.
	Domain string

	// Version is API version, e.g. v1. When creating ThirdPartyObjects will be
	// prefixed with Domain, e.g. example.com/v1.
	Version string

	// Description is free text description.
	Description string
}

// TPR allows easy operations on ThirdPartyResources. See
// https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-third-party-resource/
// for details.
type TPR struct {
	clientset kubernetes.Interface

	name          string
	group         string
	version       string
	description   string
	qualifiedName string // name.group

	endpointList string
}

func New(config Config) (*TPR, error) {
	if config.Clientset == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "k8s clientset must be set")
	}
	if config.Name == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "name must not be empty")
	}
	if config.Domain == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "group must not be empty")
	}
	if config.Version == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "version must not be empty")
	}
	if config.Description == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "description must not be empty")
	}

	tpr := &TPR{
		clientset: config.Clientset,

		name:          config.Name,
		group:         config.Domain,
		version:       config.Version,
		description:   config.Description,
		qualifiedName: config.Name + "." + config.Domain,

		endpointList: fmt.Sprintf("/apis/%s/%s/%ss", config.Domain, config.Version, config.Name),
	}
	return tpr, nil
}

// CreateAndWait create a TPR and waits till it is initialized in the cluster.
func (t *TPR) CreateAndWait() error {
	err := t.create()
	if err != nil {
		microerror.MaskAnyf(err, "creating TPR %s", t.qualifiedName)
	}
	err = t.waitInit()
	if err != nil {
		microerror.MaskAnyf(err, "waiting for TPR %s initialization", t.qualifiedName)
	}
	return nil
}

// create is extracted for testing because fake REST client does not work.
// Therefore waitInit can not be tested.
func (t *TPR) create() error {
	tpr := &v1beta1.ThirdPartyResource{
		ObjectMeta: v1.ObjectMeta{
			Name: t.qualifiedName,
		},
		Versions: []v1beta1.APIVersion{
			{Name: t.version},
		},
		Description: t.description,
	}

	_, err := t.clientset.ExtensionsV1beta1().ThirdPartyResources().Create(tpr)
	if err != nil && !errors.IsAlreadyExists(err) {
		return microerror.MaskAnyf(err, "creating TPR %s", t.qualifiedName)
	}
	return nil
}

func (t *TPR) waitInit() error {
	return util.Retry(tprInitRetryDelay, tprInitRetries, func() (bool, error) {
		_, err := t.clientset.CoreV1().RESTClient().Get().RequestURI(t.endpointList).DoRaw()
		if err != nil {
			if errors.IsNotFound(err) {
				return false, nil
			}
			return false, microerror.MaskAnyf(err, "requesting TPR %s", t.qualifiedName)
		}
		return true, nil
	})
}
