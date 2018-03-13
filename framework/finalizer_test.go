package framework

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/giantswarm/micrologger"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest/fake"
	"k8s.io/kubernetes/pkg/api/testapi"
)

func Test_addFinalizer(t *testing.T) {
	testCases := []struct {
		name           string
		object         *apiv1.Pod
		operatorName   string
		expectedResult bool
		errorMatcher   func(error) bool
	}{
		{
			name: "case 0: No finalizer is set yet",
			object: &apiv1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "TestPod",
					Namespace: "TestNamespace",
				},
			},
			operatorName:   "test-operator",
			expectedResult: true,
			errorMatcher:   nil,
		},
		{
			name: "case 1: Finalizer is already set",
			object: &apiv1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Finalizers: []string{"operatorkit.giantswarm.io/test-operator"},
					Name:       "TestPod",
					Namespace:  "TestNamespace",
				},
			},
			operatorName:   "test-operator",
			expectedResult: false,
			errorMatcher:   nil,
		},
	}
	logger, err := micrologger.New(micrologger.DefaultConfig())
	if err != nil {
		t.Fatalf("micrologger.New() failed: %#v", err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := &fake.RESTClient{
				Client: fake.CreateHTTPClient(func(req *http.Request) (*http.Response, error) {
					switch m := req.Method; {
					case m == "PATCH":
						return &http.Response{StatusCode: 200, Header: defaultHeader(), Body: objBody(testapi.Default.Codec())}, nil
					default:
						return nil, fmt.Errorf("unexpected request")
					}
				}),
				NegotiatedSerializer: serializer.DirectCodecFactory{
					CodecFactory: scheme.Codecs,
				},
			}

			f := Framework{
				logger:         logger,
				name:           tc.operatorName,
				backOffFactory: DefaultBackOffFactory(),
				restClient:     c,
			}
			result, err := f.addFinalizer(tc.object)

			switch {
			case err == nil && tc.errorMatcher == nil: // correct; carry on
			case err != nil && tc.errorMatcher != nil:
				if !tc.errorMatcher(err) {
					t.Fatalf("error == %#v, want matching", err)
				}
			case err != nil && tc.errorMatcher == nil:
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.errorMatcher != nil:
				t.Fatalf("error == nil, want non-nil")
			}

			if !reflect.DeepEqual(result, tc.expectedResult) {
				t.Fatalf("k8sEndpoint == %v, want %v", result, tc.expectedResult)
			}
		})
	}
}

func Test_removeFinalizer(t *testing.T) {
}

func defaultHeader() http.Header {
	header := http.Header{}
	header.Set("Content-Type", runtime.ContentTypeJSON)
	return header
}

func objBody(codec runtime.Codec) io.ReadCloser {
	obj := &apiv1.Pod{}
	return ioutil.NopCloser(bytes.NewReader([]byte(runtime.EncodeOrDie(codec, obj))))
}
