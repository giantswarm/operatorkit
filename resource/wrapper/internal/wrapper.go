package internal

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/resource"
)

type Wrapper interface {
	Wrapped() resource.Interface
}

func Underlying(r resource.Interface) (resource.Interface, error) {
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
