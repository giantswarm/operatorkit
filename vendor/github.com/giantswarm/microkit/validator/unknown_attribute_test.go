package validator

import (
	"testing"
)

func Test_UnknownAttribute(t *testing.T) {
	testCases := []struct {
		Received       map[string]interface{}
		Expected       map[string]interface{}
		ErrorMatcher   func(err error) bool
		ErrorAttribute string
	}{
		{
			Received:       map[string]interface{}{"foo": "foo"},
			Expected:       testExpectedForUnknownFields(t),
			ErrorAttribute: "",
			ErrorMatcher:   nil,
		},
		{
			Received:       map[string]interface{}{"foo": "baz"},
			Expected:       testExpectedForUnknownFields(t),
			ErrorAttribute: "",
			ErrorMatcher:   nil,
		},
		{
			Received:       map[string]interface{}{"foo": "foo", "bar": "baz"},
			Expected:       testExpectedForUnknownFields(t),
			ErrorAttribute: "",
			ErrorMatcher:   nil,
		},
		{
			Received:       map[string]interface{}{"foo": "foo", "three": []string{}},
			Expected:       testExpectedForUnknownFields(t),
			ErrorAttribute: "",
			ErrorMatcher:   nil,
		},
		{
			Received:       map[string]interface{}{"wrong": "field"},
			Expected:       testExpectedForUnknownFields(t),
			ErrorAttribute: "wrong",
			ErrorMatcher:   IsUnknownAttribute,
		},
		{
			Received:       map[string]interface{}{"baz": 3.54},
			Expected:       testExpectedForUnknownFields(t),
			ErrorAttribute: "baz",
			ErrorMatcher:   IsUnknownAttribute,
		},
	}

	for i, testCase := range testCases {
		err := UnknownAttribute(testCase.Received, testCase.Expected)

		if (err != nil && testCase.ErrorMatcher == nil) || (testCase.ErrorMatcher != nil && !testCase.ErrorMatcher(err)) {
			t.Fatal("case", i+1, "expected", true, "got", false)
		}

		if testCase.ErrorMatcher != nil {
			errorAttribute := ToUnknownAttribute(err).Attribute()
			if errorAttribute != testCase.ErrorAttribute {
				t.Fatal("case", i+1, "expected", testCase.ErrorAttribute, "got", errorAttribute)
			}
		}
	}
}

func Test_StructToMap(t *testing.T) {
	type testStruct struct {
		Foo string
		Bar int
	}

	s := testStruct{}

	m, err := StructToMap(s)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
	if _, ok := m["Foo"]; !ok {
		t.Fatal("expected", true, "got", false)
	}
	if _, ok := m["Bar"]; !ok {
		t.Fatal("expected", true, "got", false)
	}
}

func testExpectedForUnknownFields(t *testing.T) map[string]interface{} {
	return map[string]interface{}{
		"foo":   "foo",
		"bar":   "bar",
		"three": 3,
	}
}
