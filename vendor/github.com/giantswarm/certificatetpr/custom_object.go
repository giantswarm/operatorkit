package certificatetpr

import (
	"encoding/json"

	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/meta"
	"k8s.io/client-go/pkg/api/unversioned"
)

// CustomObject represents the Certificate TPR's custom object. It holds the
// specifications of the resource the Certificate operator is interested in.
type CustomObject struct {
	unversioned.TypeMeta `json:",inline"`
	Metadata             api.ObjectMeta `json:"metadata"`

	Spec Spec `json:"spec" yaml:"spec"`
}

// GetObjectKind is required to satisfy the Object interface.
func (c *CustomObject) GetObjectKind() unversioned.ObjectKind {
	return &c.TypeMeta
}

// GetObjectMeta is required to satisfy the ObjectMetaAccessor interface.
func (c *CustomObject) GetObjectMeta() meta.Object {
	return &c.Metadata
}

// The code below is used only to work around a known problem with third-party
// resources and ugorji. If/when these issues are resolved, the code below
// should no longer be required.

type CustomObjectCopy CustomObject

func (c *CustomObject) UnmarshalJSON(data []byte) error {
	tmp := CustomObjectCopy{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	tmp2 := CustomObject(tmp)
	*c = tmp2
	return nil
}
