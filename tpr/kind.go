package tpr

import (
	"strings"

	microerror "github.com/giantswarm/microkit/error"
)

/*
	Code based on:
	https://github.com/kubernetes/kubernetes/blob/6a4d5cd7cc58e28c20ca133dab7b0e9e56192fe3/pkg/registry/extensions/thirdpartyresourcedata/util.go
*/

// extractKindAndGroup extracts kind and group from a name. For details see the
// Config.Name godoc.
func extractKindAndGroup(name string) (kind, group string, err error) {
	parts := strings.Split(name, ".")
	if len(parts) < 3 {
		return "", "", microerror.MaskAnyf(unexpectedlyShortResourceNameError, "%s, expected at least <kind>.<domain>.<tld>", name)
	}

	// kind
	kindPart := parts[0]
	toUpper := true
	for i := range kindPart {
		char := kindPart[i]
		if toUpper {
			kind = kind + string([]byte{(char - 32)})
			toUpper = false
		} else if char == '-' {
			toUpper = true
		} else {
			kind = kind + string([]byte{char})
		}
	}

	group = strings.Join(parts[1:], ".")

	return
}

// unpluralizedSuffixes is a list of resource suffixes that are the same when
// plural and singular. This is only necessary because some bits of
// (kubernetes) code are lazy and don't actually use the RESTMapper like they
// should.
//
// Copied from:
// https://github.com/kubernetes/kubernetes/blob/b0b711119b48854e0b73805e42be2bcc4b2bd604/staging/src/k8s.io/apimachinery/pkg/api/meta/restmapper.go#L131-L137
var unpluralizedSuffixes = []string{
	"endpoints",
}

// unsafeGuessKindToResource converts Kind to a plural resource name. It is
// named after and has similar semantics to:
// https://github.com/kubernetes/kubernetes/blob/b0b711119b48854e0b73805e42be2bcc4b2bd604/staging/src/k8s.io/apimachinery/pkg/api/meta/restmapper.go#L139-L164
func unsafeGuessKindToResource(kind string) string {
	if len(kind) == 0 {
		return ""
	}

	resource := strings.ToLower(kind)

	for _, skip := range unpluralizedSuffixes {
		if strings.HasSuffix(resource, skip) {
			return resource
		}
	}

	switch string(resource[len(resource)-1]) {
	case "s":
		return resource + "es"
	case "y":
		return resource[:len(resource)-1] + "ies"
	}

	return resource + "s"
}
