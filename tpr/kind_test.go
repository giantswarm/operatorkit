package tpr

import (
	"testing"

	"github.com/juju/errgo"
	"github.com/stretchr/testify/assert"
)

func TestExtractKindAndGroup(t *testing.T) {
	tests := []struct {
		name          string
		expectedKind  string
		expectedGroup string
		expectedError error
	}{
		{
			name:          "foo.company.com",
			expectedKind:  "Foo",
			expectedGroup: "company.com",
		},
		{
			name:          "cron-tab.company.com",
			expectedKind:  "CronTab",
			expectedGroup: "company.com",
		},
		{
			name:          "foo",
			expectedError: unexpectedlyShortResourceNameError,
		},
		{
			name:          "foo.company",
			expectedError: unexpectedlyShortResourceNameError,
		},
	}

	for i, tc := range tests {
		kind, group, err := extractKindAndGroup(tc.name)
		assert.Equal(t, tc.expectedError, errgo.Cause(err), "#%d", i)
		assert.Equal(t, tc.expectedKind, kind, "#%d", i)
		assert.Equal(t, tc.expectedGroup, group, "#%d", i)
	}
}

func TestUnsafeGuessKindToResource(t *testing.T) {
	tests := []struct {
		kind             string
		expectedResource string
	}{
		{
			kind:             "Pod",
			expectedResource: "pods",
		},
		{
			kind:             "ReplicationController",
			expectedResource: "replicationcontrollers",
		},
		{
			kind:             "ImageRepository",
			expectedResource: "imagerepositories",
		},
		{
			kind:             "miss",
			expectedResource: "misses",
		},
	}

	for i, tc := range tests {
		resource := unsafeGuessKindToResource(tc.kind)
		assert.Equal(t, tc.expectedResource, resource, "#%d", i)
	}
}
