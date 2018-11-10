package admissioncontrollerendpoint

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	kitendpoint "github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"k8s.io/api/admission/v1beta1"
)

const (
	Method = "POST"
	Name   = "operatorkit/admissioncontrollerendpoint"
)

type Config struct {
	Logger   micrologger.Logger
	Reviewer Reviewer

	Path string
}

type Endpoint struct {
	logger   micrologger.Logger
	reviewer Reviewer

	path string
}

func New(config Config) (*Endpoint, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.Reviewer == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Reviewer must not be empty", config)
	}

	if config.Path == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Path must not be empty", config)
	}

	e := &Endpoint{
		logger:   config.Logger,
		reviewer: config.Reviewer,

		path: config.Path,
	}

	return e, nil
}

func (e *Endpoint) Decoder() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		var err error

		var body []byte
		{
			if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
				return nil, microerror.Maskf(decodeFailedError, "content type must be any of application/json")
			}

			body, err = ioutil.ReadAll(r.Body)
			if err != nil {
				return nil, microerror.Maskf(decodeFailedError, err.Error())
			}
		}

		var request v1beta1.AdmissionReview
		{
			_, _, err := codecs.UniversalDeserializer().Decode(body, nil, &request)
			if err != nil {
				return nil, microerror.Maskf(decodeFailedError, err.Error())
			}
		}

		return request, nil
	}
}

func (e *Endpoint) Encoder() kithttp.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		return json.NewEncoder(w).Encode(response)
	}
}

// Endpoint implements the actual logic of the admission controller HTTP
// endpoint. It expects an admission review as input that contains an admission
// request. Endpoint forwards the admission request to the configured reviewer,
// which returns an admission response. This admission response is put into a
// clean admission review and returned.
func (e *Endpoint) Endpoint() kitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var err error

		var admissionResponse *v1beta1.AdmissionResponse
		{
			admissionRequest := request.(v1beta1.AdmissionReview).Request

			admissionResponse, err = e.reviewer.Review(ctx, admissionRequest)
			if err != nil {
				return nil, microerror.Mask(err)
			}

			admissionResponse.UID = admissionRequest.UID
		}

		admissionReview := v1beta1.AdmissionReview{
			Response: admissionResponse,
		}

		return admissionReview, nil
	}
}

func (e *Endpoint) Method() string {
	return Method
}

func (e *Endpoint) Middlewares() []kitendpoint.Middleware {
	return []kitendpoint.Middleware{}
}

func (e *Endpoint) Name() string {
	return Name
}

func (e *Endpoint) Path() string {
	return e.path
}
