package name

import (
	"strconv"
	"testing"

	"github.com/giantswarm/operatorkit/controller/internal/test/handler/bar"
	"github.com/giantswarm/operatorkit/controller/internal/test/handler/foo"
	"github.com/giantswarm/operatorkit/controller/internal/test/handler/nopointer"
	"github.com/giantswarm/operatorkit/handler"
)

func Test_Handler_Name(t *testing.T) {
	testCases := []struct {
		name         string
		handler      handler.Interface
		expectedName string
	}{
		{
			name:         "case 0",
			handler:      &foo.Handler{},
			expectedName: "foo",
		},
		{
			name:         "case 1",
			handler:      &bar.Handler{},
			expectedName: "bar",
		},
		{
			name:         "case 2",
			handler:      nopointer.Handler{},
			expectedName: "nopointer",
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			name := Name(tc.handler)

			if name != tc.expectedName {
				t.Fatalf("expected handler name %#q got %#q", tc.expectedName, name)
			}
		})
	}
}
