package validator

import (
	"testing"
)

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
