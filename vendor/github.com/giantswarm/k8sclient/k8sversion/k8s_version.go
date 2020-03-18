package k8sversion

import (
	"regexp"
	"sort"
	"strconv"

	"github.com/giantswarm/microerror"
)

const (
	// The bigger the version type number, the higher the priority in the parsed
	// semver format, where alpha is major 0.0.0, beta is major 1.0.0 and GA is
	// major 2.0.0.
	versionTypeAlpha = iota
	versionTypeBeta
	versionTypeGA
)

var kubeVersionRegex = regexp.MustCompile("^v([\\d]+)(?:(alpha|beta)([\\d]+))?$")

type pair struct {
	k8s   string
	major int
	minor int
	patch int
}

// Latest returns the latest version given with the Kubernetes API versions.
// Given [v1alpha1, v1, v2alpha3] then v2alpha3 is returned as it is the latest
// version.
func Latest(versions []string) (string, error) {
	if len(versions) == 0 {
		return "", microerror.Maskf(invalidKubeVersionError, "versions must not be empty")
	}

	var pairs []pair
	for _, v := range versions {
		major, minor, patch, err := parseSemver(v)
		if err != nil {
			return "", microerror.Mask(err)
		}

		p := pair{
			k8s:   v,
			major: major,
			minor: minor,
			patch: patch,
		}

		pairs = append(pairs, p)
	}

	less := func(i, j int) bool {
		pi := pairs[i]
		pj := pairs[j]

		if pi.major < pj.major {
			return true
		} else if pi.major > pj.major {
			return false
		}

		if pi.minor < pj.minor {
			return true
		} else if pi.minor > pj.minor {
			return false
		}

		if pi.patch < pj.patch {
			return true
		} else if pi.patch < pj.patch {
			return false
		}

		return false
	}

	sort.Slice(pairs, less)

	return pairs[len(pairs)-1].k8s, nil
}

// Less returns true if a is less/older that b, where a and b are Kubernetes API
// Versions. Given v1 and v2alpha3 then v1 is returned since it is the
// semantically older version.
func Less(a string, b string) (bool, error) {
	if a == b {
		return false, nil
	}

	latest, err := Latest([]string{a, b})
	if err != nil {
		return false, microerror.Mask(err)
	}

	return b == latest, nil
}

// parseSemver takes a version string which is meant to express a Kubernetes
// APIVersion. The implementation is heavily inspired by the apimachinery
// repository upstream provides. For more information on Kubernetes API
// Versioning check the official docs.
//
//     https://kubernetes.io/docs/concepts/overview/kubernetes-api/#api-versioning
//     https://github.com/kubernetes/apimachinery/blob/b9f0d37e94c6953b55a668a9c4134da6262acfe5/pkg/version/helpers.go
//
func parseSemver(v string) (major int, minor int, patch int, err error) {
	submatches := kubeVersionRegex.FindStringSubmatch(v)
	if len(submatches) != 4 {
		return 0, 0, 0, microerror.Maskf(invalidKubeVersionError, v)
	}

	switch submatches[2] {
	case "alpha":
		minor = versionTypeAlpha
	case "beta":
		minor = versionTypeBeta
	case "":
		minor = versionTypeGA
	default:
		return 0, 0, 0, microerror.Maskf(invalidKubeVersionError, v)
	}

	major, err = strconv.Atoi(submatches[1])
	if err != nil {
		return 0, 0, 0, microerror.Maskf(invalidKubeVersionError, v)
	}
	if minor != versionTypeGA {
		patch, err = strconv.Atoi(submatches[3])
		if err != nil {
			return 0, 0, 0, microerror.Maskf(invalidKubeVersionError, v)
		}
	}

	return major, minor, patch, nil
}
