package validator

import (
	"fmt"
	"sort"

	microerror "github.com/giantswarm/microkit/error"
	"github.com/juju/errgo"
)

// ImmutableAttributeError indicates a data structure is invalid because it
// contains fields that are immutable.
type ImmutableAttributeError struct {
	attribute string
	message   string
}

// Attribute returns the attribute that is causing the immutable attribute error.
func (e ImmutableAttributeError) Attribute() string {
	return e.attribute
}

// Error returns ImmutableAttributeError's message.
// This way ImmutableAttributeError implements the error interface.
func (e ImmutableAttributeError) Error() string {
	return e.message
}

// IsImmutableAttributeError lets you check if an error is an ImmutableAttributeError.
func IsImmutableAttributeError(err error) bool {
	_, ok := errgo.Cause(err).(ImmutableAttributeError)
	return ok
}

// ToImmutableAttributeError tries to cast a given error into a ImmutableAttributeError and
// returns it. ToImmutableAttributeError will panic in case the underlying error is not of
// type ToImmutableAttributeError. Use IsImmutableAttributeError before calling ToImmutableAttributeError.
func ToImmutableAttributeError(err error) ImmutableAttributeError {
	return errgo.Cause(err).(ImmutableAttributeError)
}

// ValidateImmutableAttribute takes an arbitrary map and a map obtaining some expected
// structure. The first argument might represent an incoming request of some
// microservice. The second argument should contain a datastructure representing
// only the attributes that are allowed to be mutated. If the first map contains
// fields which are not in the whitelist expected, an ImmutableAttributeError
// is returned.
func ValidateImmutableAttribute(received, blacklist map[string]interface{}) error {
	alphabeticalBlacklist := make([]string, 0, len(blacklist))

	for key := range blacklist {
		alphabeticalBlacklist = append(alphabeticalBlacklist, key)
	}

	sort.Strings(alphabeticalBlacklist)

	for _, blacklistedKey := range alphabeticalBlacklist {
		var found bool

		for recievedKey := range received {
			if recievedKey == blacklistedKey {
				found = true
				break
			}
		}

		if found {
			err := ImmutableAttributeError{
				attribute: blacklistedKey,
				message:   fmt.Sprintf("attribute '%s' is immutable, you are not allowed to change it", blacklistedKey),
			}

			return microerror.MaskAny(err)
		}

	}

	return nil
}
