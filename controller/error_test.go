package controller

import (
	"testing"

	"github.com/giantswarm/microerror"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var unknownError = microerror.New("unknown error")

func Test_IsStatusForbidden(t *testing.T) {
	testCases := []struct {
		name           string
		err            error
		expectedResult bool
	}{
		{
			name:           "case 0: match apimachinery StatusError with Forbidden reason",
			err:            errors.NewForbidden(schema.GroupResource{}, "unittest", nil),
			expectedResult: true,
		},
		{
			name:           "case 1: match masked apimachinery StatusError with Forbidden reason",
			err:            microerror.Mask(errors.NewForbidden(schema.GroupResource{}, "unittest", nil)),
			expectedResult: true,
		},
		{
			name:           "case 2: don't match nil",
			err:            nil,
			expectedResult: false,
		},
		{
			name:           "case 3: don't match unknown error",
			err:            unknownError,
			expectedResult: false,
		},
		{
			name:           "case 4: don't match apimachinery StatusError with Unauthorized reason",
			err:            errors.NewUnauthorized(""),
			expectedResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := IsStatusForbidden(tc.err)

			if result != tc.expectedResult {
				t.Fatalf("IsStatusForbidden(%#v) == %v, expected %v", tc.err, result, tc.expectedResult)
			}
		})
	}
}
