package internal

import (
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/operatorkit/controller"
)

type Wrapper interface {
	Wrapped() controller.Resource
}

func Underlying(r controller.Resource) (controller.Resource, error) {
	i := 0

	for {
		wrapper, ok := r.(Wrapper)
		if ok {
			r = wrapper.Wrapped()
		} else {
			return r, nil
		}

		// When more that 1000 interations, assume infifinite loop.
		i++
		if i > 1000 {
			return nil, microerror.Maskf(loopDetectedError, "unwrapped 1000 times, assuming infite loop: resource = %s", r.Name())
		}
	}
}
