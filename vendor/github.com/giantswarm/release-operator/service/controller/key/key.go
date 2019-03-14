package key

import (
	"strings"

	applicationv1alpha1 "github.com/giantswarm/apiextensions/pkg/apis/application/v1alpha1"
	releasev1alpha1 "github.com/giantswarm/apiextensions/pkg/apis/release/v1alpha1"
	"github.com/giantswarm/microerror"
)

const (
	// Namespace is the namespace where App CRs are created.
	Namespace = "giantswarm"

	LabelAppOperatorVersion = "app-operator.giantswarm.io/version"
	LabelManagedBy          = "giantswarm.io/managed-by"
	LabelReleaseCyclePhase  = "release-operator.giantswarm.io/release-cycle-phase"
	LabelServiceType        = "giantswarm.io/service-type"

	ValueServiceTypeManaged = "managed"
)

// ReleaseVersion returns the version of the given release.
func ReleaseVersion(releaseCR releasev1alpha1.Release) string {
	return releaseCR.Spec.Version
}

// SplitReleaseName splits a release name into provider and version.
// It returns provider, version, and error, in this order.
//
// It expects name to be in the following format <provider>.<version>
// e.g. aws.v6.0.1
func SplitReleaseName(name string) (string, string, error) {
	split := strings.SplitN(name, ".", 2)
	if len(split) < 2 || len(split[0]) == 0 || len(split[1]) == 0 {
		return "", "", microerror.Maskf(invalidReleaseNameError, "expect <provider>.<version>, got %#q", name)
	}

	return split[0], split[1], nil
}

// ToAppCR converts v into an App CR.
func ToAppCR(v interface{}) (*applicationv1alpha1.App, error) {
	x, ok := v.(*applicationv1alpha1.App)
	if !ok {
		return nil, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", x, v)
	}

	return x, nil
}

// ToReleaseCR converts v into a Release CR.
func ToReleaseCR(v interface{}) (*releasev1alpha1.Release, error) {
	x, ok := v.(*releasev1alpha1.Release)
	if !ok {
		return nil, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", x, v)
	}

	return x, nil
}

// ToReleaseCycleCR converts v into a ReleaseCycle CR.
func ToReleaseCycleCR(v interface{}) (*releasev1alpha1.ReleaseCycle, error) {
	x, ok := v.(*releasev1alpha1.ReleaseCycle)
	if !ok {
		return nil, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", x, v)
	}

	return x, nil
}
