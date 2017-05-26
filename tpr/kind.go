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
