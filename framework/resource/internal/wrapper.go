package internal

import "github.com/giantswarm/operatorkit/framework"

type Wrapper interface {
	Underlying() framework.Resource
}
