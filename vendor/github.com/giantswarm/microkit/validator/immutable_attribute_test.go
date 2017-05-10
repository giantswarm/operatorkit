package validator

import (
	"testing"
)

func Test_ValidateImmutableAttribute(t *testing.T) {
	testCases := []struct {
		Received       map[string]interface{}
		Blacklist      map[string]interface{}
		ErrorMatcher   func(err error) bool
		ErrorAttribute string
	}{
		// With no black list provided, it shouldn't complain at all.
		{
			Received:       map[string]interface{}{"name": "Oliver", "age": 29},
			Blacklist:      map[string]interface{}{},
			ErrorAttribute: "",
			ErrorMatcher:   nil,
		},

		// With an attribute in the black list that is not present in the received
		// struct, it shouldn't complain at all.
		{
			Received:       map[string]interface{}{"name": "Oliver", "age": 29},
			Blacklist:      map[string]interface{}{"hairColor": ""},
			ErrorAttribute: "",
			ErrorMatcher:   nil,
		},

		// With an attribute in the black list that _is_ present in the received
		// struct, it should return a ImmutableAttribute error.
		{
			Received:       map[string]interface{}{"name": "Oliver", "age": 29},
			Blacklist:      map[string]interface{}{"name": ""},
			ErrorAttribute: "name",
			ErrorMatcher:   IsImmutableAttributeError,
		},

		// With multiple attributes in the black list that are present in the received
		// struct, it should return a ImmutableAttribute error for the first black listed
		// attribute ordered alphabetically.
		{
			Received:       map[string]interface{}{"name": "Oliver", "age": 29, "zebra": true},
			Blacklist:      map[string]interface{}{"name": "", "age": 0, "zebra": false},
			ErrorAttribute: "age",
			ErrorMatcher:   IsImmutableAttributeError,
		},
	}

	for i, testCase := range testCases {
		err := ValidateImmutableAttribute(testCase.Received, testCase.Blacklist)

		if (err != nil && testCase.ErrorMatcher == nil) || (testCase.ErrorMatcher != nil && !testCase.ErrorMatcher(err)) {
			t.Fatal("case", i+1, "expected", true, "got", false)
		}

		if testCase.ErrorMatcher != nil {
			errorAttribute := ToImmutableAttributeError(err).Attribute()
			if errorAttribute != testCase.ErrorAttribute {
				t.Fatal("case", i+1, "expected", testCase.ErrorAttribute, "got", errorAttribute)
			}
		}
	}
}
