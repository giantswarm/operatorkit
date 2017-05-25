package tpr

import (
	"testing"

	"github.com/juju/errgo"
	"github.com/stretchr/testify/assert"
)

func TestExtractKindAndGroup(t *testing.T) {
	tests := []struct {
		Input         string
		ExpectedKind  string
		ExpectedGroup string
		ExpectedError error
	}{
		{
			Input:         "foo.company.com",
			ExpectedKind:  "Foo",
			ExpectedGroup: "company.com",
		},
		{
			Input:         "cron-tab.company.com",
			ExpectedKind:  "CronTab",
			ExpectedGroup: "company.com",
		},
		{
			Input:         "foo",
			ExpectedError: unexpectedlyShortResourceNameError,
		},
		{
			Input:         "foo.company",
			ExpectedError: unexpectedlyShortResourceNameError,
		},
	}

	for i, tt := range tests {
		kind, group, err := extractKindAndGroup(tt.Input)
		assert.Equal(t, tt.ExpectedError, errgo.Cause(err), "#%d", i)
		assert.Equal(t, tt.ExpectedKind, kind, "#%d", i)
		assert.Equal(t, tt.ExpectedGroup, group, "#%d", i)
	}
}
