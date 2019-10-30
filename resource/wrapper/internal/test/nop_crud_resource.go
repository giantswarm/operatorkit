package test

import (
	"fmt"

	"github.com/giantswarm/micrologger/microloggertest"

	"github.com/giantswarm/operatorkit/resource"
	"github.com/giantswarm/operatorkit/resource/crud"
)

func NewNopCRUDResource() resource.Interface {
	c := crud.ResourceConfig{
		CRUD:   NewNopCRUD(),
		Logger: microloggertest.New(),
	}

	r, err := crud.NewResource(c)
	if err != nil {
		panic(fmt.Sprintf("%#v", err))
	}

	return r
}
