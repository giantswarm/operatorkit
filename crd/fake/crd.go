package fake

import (
	"fmt"

	"github.com/giantswarm/operatorkit/crd"
)

func MustNewCRD() *crd.CRD {
	crdConfig := crd.DefaultConfig()

	crdConfig.Group = Group
	crdConfig.Kind = Kind
	crdConfig.Name = Name
	crdConfig.Plural = Plural
	crdConfig.Singular = Singular
	crdConfig.Scope = Scope
	crdConfig.Version = VersionV1

	newCRD, err := crd.New(crdConfig)
	if err != nil {
		panic(fmt.Sprintf("%#v", err))
	}

	return newCRD
}
