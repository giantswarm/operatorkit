package validator

import (
	"encoding/json"

	microerror "github.com/giantswarm/microkit/error"
)

// StructToMap is a helper method to convert an expected request data structure
// in the correctly formatted type to UnknownAttributes.
func StructToMap(s interface{}) (map[string]interface{}, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	var m map[string]interface{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	return m, nil
}
