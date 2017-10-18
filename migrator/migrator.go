package migrator

import (
	"context"
	"encoding/json"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
)

// TransformFunc transforms Kubernetes objects from one version to another. It
// is guaranteed that all the objects passed to the function have the same
// group, version, and kind. TransformFunc is required to return objects with
// the same group and kind. All objects returned by TransfomFunc must have the
// same version.
type TransformFunc func(context.Context, []runtime.Object) []runtime.Object

type Config struct {
	// Dependencies.

	K8sClient rest.Interface
	// NewActualZeroObjectFunc creates zero-value object to be transformed
	// by the TransformFunc.
	NewActaulZeroObjectFunc func() runtime.Object
	// TransformFunc is the transformation function applied on objects
	// which are subject of the migration. See TransofrmFunc godoc for
	// details.
	TransformFunc TransformFunc

	// Settings.

	// Endpoint is the endpoint of the resource to be migrated.
	Endpoint string
	// ActualVersion is the actual version of the resource to be migrated.
	ActualVersion string
	// DesiredVersion is the version of the resource after transformation.
	DesiredVersion string
}

func DefaultConfig() Config {
	// TODO default config
}

type Migrator struct {
	// Dependencies.

	logger                  micrologger.Logger
	k8sClient               rest.Interface
	transformFunc           TransformFunc
	newActualZeroObjectFunc func() runtime.Object

	// Settings.

	endpoint string
}

func New(config Config) (*Migrator, error) {
	// TODO invalidConfig

	logger := config.Logger.With("component", "Migrator")
}

func (m *Migrator) Migrate(ctx context.Context) error {
	m.logger.Log("function", "Migrate", "action", "start")

	var current []runtime.Object
	{
		bytes, err := m.k8sClient.
			Get().
			Context(ctx).
			AbsPath(m.endpoint).
			DoRaw()

		if err != nil {
			return nil, microerror.Maskf(err, "get TPO")
		}
		if ctx.Err() != nil {
			return microerror.Mask(ctx.Err())
		}

		v := m.newActualZeroObjectFunc()
		err = json.Unmarshal(bytes, &v)
		if err != nil {
			return nil, microerror.Maskf(err, "unmarshal TPO")
		}

		current = append(current, v)
	}

	if len(current) == 0 {
		m.logger.Log("function", "Migrate", "action", "end", "debug", "migrated 0 objects")
		return nil
	}

	var group, kind string
	{
		var version string
		gvk := current[0].GetObjectKind().GroupVersionKind()
		group, version, kind = gvk.Group, gvk.Version, gvk.Kind

		if version != m.actualVersion {
			return microerror.Maskf(transformationError,
				"objects to transform expected version is %q, got %q", m.actualVersion, version)
		}
	}

	transformed, err := m.transformFunc(ctx, objects)
	if err != nil {
		return microerror.Mask(err)
	}
	if ctx.Err() != nil {
		return microerror.Mask(ctx.Err())
	}
	for _, t := range transformed {
		gvk := t.GetObjectKind().GroupVersionKind()
		if group != gvk.Group {
			return microerror.Maskf(transformationError,
				"transformed objects expected group is %q, got %q", group, gvk.Group)
		}
		if version != gvk.Version {
			return microerror.Maskf(transformationError,
				"transformed objects expected version is %q, got %q", m.desiredVersion, gvk.Version)
		}
		if kind != gvk.Kind {
			return microerror.Maskf(transformationError,
				"transformed objects expected kind is %q, got %q", kind, gvk.Kind)
		}
	}

	// TODO test migration of already migrated object
	// TODO test migration with conflicting object
	//m.k8sClient.Post().
}
