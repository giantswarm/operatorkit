package internal

import (
	"testing"

	"github.com/giantswarm/operatorkit/resource/wrapper/internal/test"
)

func Test_CRUD_success(t *testing.T) {
	// For wrapping resources to work correctly CRUD must be able to
	// successfully extract crud.Interface from *crud.Resource.

	r := test.NewNopCRUDResource()

	extractedCRUD, ok := CRUD(r)
	if !ok {
		t.Fatalf("CURD(r) == %v, want %v", ok, true)
	}
	if extractedCRUD.Name() != r.Name() {
		t.Fatalf("extractedCRUD.Name() == %v, want %v", extractedCRUD.Name(), r.Name())
	}
}

func Test_CRUD_failure(t *testing.T) {
	r := test.NewNopBasicResource()

	_, ok := CRUD(r)
	if ok {
		t.Fatalf("CURD(r) == %v, want %v", ok, false)
	}
}
