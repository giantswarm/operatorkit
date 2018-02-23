package internal

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/framework"
)

type Wrapper interface {
	Wrapped() framework.Resource
}

func Underlying(r framework.Resource) (framework.Resource, error) {
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
