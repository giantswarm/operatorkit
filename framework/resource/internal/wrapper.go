package internal

import (
	old "github.com/giantswarm/operatorkit/framework"
	"github.com/giantswarm/operatorkit/framework/resource/internal/framework"
)

type OldWrapper interface {
	Underlying() old.Resource
}

func OldUnderlying(r old.Resource) old.Resource {
	for {
		wrapper, ok := r.(OldWrapper)
		if ok {
			r = wrapper.Underlying()
		} else {
			return r
		}
	}
}

type Wrapper interface {
	Wrapped() framework.Resource
}

func Underlying(r framework.Resource) framework.Resource {
	for {
		wrapper, ok := r.(Wrapper)
		if ok {
			r = wrapper.Wrapped()
		} else {
			return r
		}
	}
}
