package v1alpha1

import (
	"time"
)

// DeepCopyDate is a date type designed to be validated with json-schema date
// type.
type DeepCopyDate struct {
	time.Time
}

// MarshalJSON implements the json.Marshaler interface. The time is
// expected to be a quoted string in yyyy-mm-dd format.
//
// NOTE: This method has a value (not pointer) receiver. Otherwise marshalling
// will stop working for values. When this is a value receiver it works for both.
func (d DeepCopyDate) MarshalJSON() ([]byte, error) {
	s := d.Format(`"2006-01-02"`)
	return []byte(s), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface. The time is
// expected to be a quoted string in yyyy-mm-dd format.
func (d *DeepCopyDate) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}

	// Error masking is skipped here as this will go trough generated
	// unmarshaling code.
	var err error
	d.Time, err = time.Parse(`"2006-01-02"`, string(data))
	return err
}

func (d *DeepCopyDate) DeepCopyInto(out *DeepCopyDate) {
	*out = *d
}
