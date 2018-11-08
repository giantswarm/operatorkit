package microerror

import "encoding/json"

// Error is a predefined error structure whose purpose is to act as container
// for meta information associated to a specific error. The specific error type
// matching can be used as usual. The usual error masking and cause gathering
// can be used as usual. Using Error might look as follows. In the beginning is
// a usual error defined, along with its matcher. This error is the root cause
// once emitted during runtime.
//
//     var notEnoughWorkersError = &microerror.Error{
//         Desc: "The amount of requested tenant cluster workers exceeds the available number of control plane nodes.",
//         Docs: "https://github.com/giantswarm/ops-recipes/blob/master/349-not-enough-workers.md",
//         Kind: "notEnoughWorkersError",
//     }
//
//     func IsNotEnoughWorkers(err error) bool {
//         return microerror.Cause(err) == notEnoughWorkersError
//     }
//
type Error struct {
	Desc string `json:"desc"`
	Docs string `json:"docs"`
	Kind string `json:"kind"`
}

func (e *Error) Error() string {
	return toStringCase(e.Kind)
}

func (e *Error) GoString() string {
	return e.String()
}

func (e *Error) String() string {
	b, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}

	return string(b)
}
