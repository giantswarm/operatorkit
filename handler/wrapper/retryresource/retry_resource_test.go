package retryresource

import (
	"testing"

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/micrologger/microloggertest"

	"github.com/giantswarm/operatorkit/resource/wrapper/internal"
	"github.com/giantswarm/operatorkit/resource/wrapper/internal/test"
)

// Test_CRUD_success tests if wrapping CRUD resource allows extracting
// crud.Interface from the wrapping resource.
func Test_CRUD_success(t *testing.T) {
	var err error

	r := test.NewNopCRUDResource()

	c := Config{
		BackOff:  backoff.NewConstant(1, 1),
		Logger:   microloggertest.New(),
		Resource: r,
	}
	wrapped, err := New(c)
	if err != nil {
		t.Fatalf("err = %#v, want nil", err)
	}

	extractedCRUD, ok := internal.CRUD(wrapped)
	if !ok {
		t.Fatalf("CURD(r) == %v, want %v", ok, true)
	}
	if extractedCRUD.Name() != r.Name() {
		t.Fatalf("extractedCRUD.Name() == %v, want %v", extractedCRUD.Name(), r.Name())
	}
}

// Test_CRUD_failure tests if wrapping basic resource does not allow extracting
// crud.Interface from the wrapping resource.
func Test_CRUD_failure(t *testing.T) {
	var err error

	r := test.NewNopBasicResource()

	c := Config{
		BackOff:  backoff.NewConstant(1, 1),
		Logger:   microloggertest.New(),
		Resource: r,
	}
	wrapped, err := New(c)
	if err != nil {
		t.Fatalf("err = %#v, want nil", err)
	}

	_, ok := internal.CRUD(wrapped)
	if ok {
		t.Fatalf("Basic(r) == %v, want %v", ok, false)
	}
}
