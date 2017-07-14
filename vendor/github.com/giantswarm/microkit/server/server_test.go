package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	kitendpoint "github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
)

// Test_Server_Endpoints ensures the endpoint registration works as expected.
// There have been issues with registering endpoints, where one endpoint
// overwrote the other. This test makes sure this does not happen again. So we
// register two endpoints and check if the endpoints are actually called as
// expected. When they return the correct response, the registration is
// considered correct.
func Test_Server_Endpoints(t *testing.T) {
	e1 := testNewEndpoint(t)
	e1.(*testEndpoint).endpointResponseFormat = "e1-test-response-%d"
	e1.(*testEndpoint).path = "/e1-test-path"
	e2 := testNewEndpoint(t)
	e2.(*testEndpoint).endpointResponseFormat = "e2-test-response-%d"
	e2.(*testEndpoint).path = "/e2-test-path"

	config := DefaultConfig()
	config.Endpoints = []Endpoint{e1, e2}
	newServer, err := New(config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	newServer.Boot()
	defer newServer.Shutdown()

	{
		r, err := http.NewRequest("GET", "/e1-test-path", nil)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
		w := httptest.NewRecorder()

		newServer.Config().Router.ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Fatal("expected", http.StatusOK, "got", w.Code)
		}

		if w.Body.String() != "e1-test-response-1" {
			t.Fatal("expected", "e1-test-response-1", "got", w.Body.String())
		}
	}

	{
		r, err := http.NewRequest("GET", "/e2-test-path", nil)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
		w := httptest.NewRecorder()

		newServer.Config().Router.ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Fatal("expected", http.StatusOK, "got", w.Code)
		}

		if w.Body.String() != "e2-test-response-1" {
			t.Fatal("expected", "e2-test-response-1", "got", w.Body.String())
		}
	}
}

// Test_Server_Default_HandlerWrapper verifies the default HandlerWrapper does
// not do anything. This is the negative test for Test_Server_Custom_HandlerWrapper.
func Test_Server_Default_HandlerWrapper(t *testing.T) {
	e := testNewEndpoint(t)

	config := DefaultConfig()
	config.Endpoints = []Endpoint{e}
	newServer, err := New(config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	newServer.Boot()
	defer newServer.Shutdown()

	{
		r, err := http.NewRequest("GET", "/test-path", nil)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
		w := httptest.NewRecorder()

		newServer.Config().Router.ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Fatal("expected", http.StatusOK, "got", w.Code)
		}

		h := w.HeaderMap.Get("X-Test-Header")
		if h != "" {
			t.Fatal("expected", "no header", "got", h)
		}
	}
}

// Test_Server_Custom_HandlerWrapper verifies that the custom HandlerWrapper
// does what it is supposed to do. In this test case it sets an additional
// header to the request. In case our response recorder contains our expected
// header, the test succeeds.
func Test_Server_Custom_HandlerWrapper(t *testing.T) {
	e := testNewEndpoint(t)

	config := DefaultConfig()
	config.Endpoints = []Endpoint{e}
	config.HandlerWrapper = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test-Header", "test value")
			h.ServeHTTP(w, r)
		})
	}
	newServer, err := New(config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	newServer.Boot()
	defer newServer.Shutdown()

	{
		r, err := http.NewRequest("GET", "/test-path", nil)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
		w := httptest.NewRecorder()

		newServer.Config().Router.ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Fatal("expected", http.StatusOK, "got", w.Code)
		}

		h := w.HeaderMap.Get("X-Test-Header")
		if h != "test value" {
			t.Fatal("expected", "no header", "got", h)
		}
	}
}

type testEndpoint struct {
	decoderExecuted        int
	decoderRequest         string
	endpointExecuted       int
	encoderExecuted        int
	endpointResponseFormat string
	method                 string
	name                   string
	path                   string
}

func testNewEndpoint(t *testing.T) Endpoint {
	newEndpoint := &testEndpoint{
		decoderExecuted:        0,
		decoderRequest:         "",
		endpointExecuted:       0,
		encoderExecuted:        0,
		endpointResponseFormat: "test-response-%d",
		method:                 "GET",
		name:                   "test-endpoint",
		path:                   "/test-path",
	}

	return newEndpoint
}

func (e *testEndpoint) Decoder() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		e.decoderExecuted++
		return e.decoderRequest, nil
	}
}

func (e *testEndpoint) Encoder() kithttp.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		e.encoderExecuted++
		_, err := w.Write([]byte(response.(string)))
		return err
	}
}

func (e *testEndpoint) Endpoint() kitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		e.endpointExecuted++
		return fmt.Sprintf(e.endpointResponseFormat, e.endpointExecuted), nil
	}
}

func (e *testEndpoint) Method() string {
	return e.method
}

func (e *testEndpoint) Middlewares() []kitendpoint.Middleware {
	return []kitendpoint.Middleware{}
}

func (e *testEndpoint) Name() string {
	return e.name
}

func (e *testEndpoint) Path() string {
	return e.path
}
