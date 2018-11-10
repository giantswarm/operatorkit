package admissioncontrollerendpoint

import (
	"context"

	"k8s.io/api/admission/v1beta1"
)

type Reviewer interface {
	Review(ctx context.Context, ar *v1beta1.AdmissionRequest) (*v1beta1.AdmissionResponse, error)
}
