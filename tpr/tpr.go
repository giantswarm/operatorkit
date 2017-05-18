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
	Name        string
	Group       string
	Version     string
	Description string
}

type TPR struct {
	clientset kubernetes.Interface

	name          string
	group         string
	version       string
	description   string
	qualifiedName string

	endpointList string
}

func New(config Config) (*TPR, error) {
	if config.Clientset == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "k8s clientset must be set")
	}
	if config.Name == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "name must not be empty")
	}
	if config.Group == "" {
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
		group:         config.Group,
		version:       config.Version,
		description:   config.Description,
		qualifiedName: config.Name + "." + config.Group,

		endpointList: fmt.Sprintf("/apis/%s/%s/%ss", config.Group, config.Version, config.Name),
	}
	return tpr, nil
}

// CreateAndWait create a TPR and waits till it is initialized in the cluster.
func (t *TPR) CreateAndWait() error {
	err := t.create()
	if err != nil {
		microerror.MaskAny(fmt.Errorf("creating TPR: %+v", err))
	}
	err = t.waitInit()
	if err != nil {
		microerror.MaskAny(fmt.Errorf("waiting TPR initialization: %+v", err))
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
		return microerror.MaskAnyf(err, "creating TPR %s", t.name)
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
			return false, microerror.MaskAnyf(err, "requesting TPR %s", t.name)
		}
		return true, nil
	})
}
