package internal

import (
	"github.com/giantswarm/operatorkit/v2/resource"
	"github.com/giantswarm/operatorkit/v2/resource/crud"
)

func CRUD(r resource.Interface) (crud.Interface, bool) {
	type cruder interface {
		CRUD() crud.Interface
	}

	c, ok := r.(cruder)
	if ok {
		return c.CRUD(), true
	}

	return nil, false
}
