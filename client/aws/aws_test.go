package aws

import (
	"fmt"
	"testing"

	"github.com/giantswarm/microerror"
	"github.com/stretchr/testify/assert"
)

func TestValidAmazonAccountID(t *testing.T) {
	tests := []struct {
		name            string
		amazonAccountID string
		err             error
	}{
		{
			name:            "ID has wrong length",
			amazonAccountID: "foo",
			err:             wrongAmazonAccountIDLengthError,
		},
		{
			name:            "ID contains letters",
			amazonAccountID: "123foo123foo",
			err:             malformedAmazonAccountIDError,
		},
		{
			name:            "ID is empty",
			amazonAccountID: "",
			err:             emptyAmazonAccountIDError,
		},
		{
			name:            "ID has correct format",
			amazonAccountID: "123456789012",
			err:             nil,
		},
	}

	for _, tc := range tests {
		err := validateAccountID(tc.amazonAccountID)
		assert.Equal(t, microerror.Cause(tc.err), microerror.Cause(err), fmt.Sprintf("[%s] The return value was not what we expected", tc.name))
	}
}
