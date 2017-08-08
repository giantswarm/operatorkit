package metricsresource

import (
	"testing"
)

// Test_MetricsResource_toCamelCase ensures the resource's Prometheus related
// configurations are properly formatted as expected.
func Test_MetricsResource_toCamelCase(t *testing.T) {
	testCases := []struct {
		InputString    string
		ExpectedString string
	}{
		{
			InputString:    "foo-bar",
			ExpectedString: "fooBar",
		},
		{
			InputString:    "foo+bar",
			ExpectedString: "fooBar",
		},
		{
			InputString:    "foo bar",
			ExpectedString: "fooBar",
		},
		{
			InputString:    "foo_bar",
			ExpectedString: "fooBar",
		},
		{
			InputString:    "fooBar",
			ExpectedString: "fooBar",
		},
	}

	for i, tc := range testCases {
		output := toCamelCase(tc.InputString)
		if output != tc.ExpectedString {
			t.Fatal("test", i+1, "expected", tc.ExpectedString, "got", output)
		}
	}
}
