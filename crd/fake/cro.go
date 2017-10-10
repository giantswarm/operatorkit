package fake

import (
	"encoding/json"
	"fmt"

	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewCRO(ID string) []byte {
	newCustomObject := &CustomObject{
		TypeMeta: apismetav1.TypeMeta{
			APIVersion: Group + `/` + VersionV1,
			Kind:       Kind,
		},
		ObjectMeta: apismetav1.ObjectMeta{
			Name: ID,
		},
		Spec: Spec{
			ID: ID,
		},
	}

	b, err := json.Marshal(newCustomObject)
	if err != nil {
		panic(fmt.Sprintf("%#v", err))
	}

	return b
}
