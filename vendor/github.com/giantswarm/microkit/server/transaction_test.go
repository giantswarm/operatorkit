package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/net/context"
)

func Test_Transaction_IDFormat(t *testing.T) {
	testCases := []struct {
		TransactionID string
		Valid         bool
	}{
		{
			TransactionID: "",
			Valid:         false,
		},
		{
			TransactionID: "foo",
			Valid:         false,
		},
		{
			TransactionID: "d.99ab4af-ddc7-4c7b-8e2b-1cdef5b129c7",
			Valid:         false,
		},
		{
			TransactionID: "-99ab4af-ddc7-4c7b-8e2b-1cdef5b129c7",
			Valid:         false,
		},
		{
			TransactionID: "d99ab4af-ddc7-4c7b-8e2b-1cdef5b129c-",
			Valid:         false,
		},
		{
			TransactionID: "d99ab4af-ddc7-4c7b-8e2b-1cdef5b129c7",
			Valid:         true,
		},
		{
			TransactionID: "a1e0d43b-fea2-4240-84a7-7abdffca1999",
			Valid:         true,
		},
		{
			TransactionID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			Valid:         true,
		},
		{
			TransactionID: "00000000-0000-0000-0000-000000000000",
			Valid:         true,
		},
	}

	for _, testCase := range testCases {
		isValid := IsValidTransactionID(testCase.TransactionID)
		if isValid != testCase.Valid {
			t.Fatal("expected", testCase.Valid, "got", isValid)
		}
	}
}

func Test_Transaction_NoIDGiven(t *testing.T) {
	e := testNewEndpoint(t)

	config := DefaultConfig()
	config.Endpoints = []Endpoint{e}
	newServer, err := New(config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	newServer.Boot()
	defer newServer.Shutdown()

	// Here we make a request against our test endpoint. The endpoint is executed
	// the first time. So the execution counts should be one.
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

		decoderExecuted := e.(*testEndpoint).decoderExecuted
		if decoderExecuted != 1 {
			t.Fatal("expected", 1, "got", decoderExecuted)
		}
		endpointExecuted := e.(*testEndpoint).endpointExecuted
		if endpointExecuted != 1 {
			t.Fatal("expected", 1, "got", endpointExecuted)
		}
		encoderExecuted := e.(*testEndpoint).encoderExecuted
		if encoderExecuted != 1 {
			t.Fatal("expected", 1, "got", encoderExecuted)
		}

		if w.Body.String() != "test-response-1" {
			t.Fatal("expected", "test-response-1", "got", w.Body.String())
		}
	}

	// Here we make another request against our test endpoint. The endpoint is
	// executed the second time. So the execution counts should be two.
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

		decoderExecuted := e.(*testEndpoint).decoderExecuted
		if decoderExecuted != 2 {
			t.Fatal("expected", 2, "got", decoderExecuted)
		}
		endpointExecuted := e.(*testEndpoint).endpointExecuted
		if endpointExecuted != 2 {
			t.Fatal("expected", 2, "got", endpointExecuted)
		}
		encoderExecuted := e.(*testEndpoint).encoderExecuted
		if encoderExecuted != 2 {
			t.Fatal("expected", 2, "got", encoderExecuted)
		}

		if w.Body.String() != "test-response-2" {
			t.Fatal("expected", "test-response-2", "got", w.Body.String())
		}
	}
}

func Test_Transaction_IDGiven(t *testing.T) {
	e := testNewEndpoint(t)

	config := DefaultConfig()
	config.Endpoints = []Endpoint{e}
	newServer, err := New(config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	newServer.Boot()
	defer newServer.Shutdown()

	// Here we make a request against our test endpoint. The endpoint is executed
	// the first time. So the execution counts should be one.
	{
		r, err := http.NewRequest("GET", "/test-path", nil)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
		r.Header.Add(TransactionIDHeader, "my-very-valid-test-transaction-id")
		w := httptest.NewRecorder()

		newServer.Config().Router.ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Fatal("expected", http.StatusOK, "got", w.Code)
		}

		decoderExecuted := e.(*testEndpoint).decoderExecuted
		if decoderExecuted != 1 {
			t.Fatal("expected", 1, "got", decoderExecuted)
		}
		endpointExecuted := e.(*testEndpoint).endpointExecuted
		if endpointExecuted != 1 {
			t.Fatal("expected", 1, "got", endpointExecuted)
		}
		encoderExecuted := e.(*testEndpoint).encoderExecuted
		if encoderExecuted != 1 {
			t.Fatal("expected", 1, "got", encoderExecuted)
		}

		if w.Body.String() != "test-response-1" {
			t.Fatal("expected", "test-response-1", "got", w.Body.String())
		}
	}

	// Here we make another request against our test endpoint. In this request and
	// the previous one we provided the same transaction ID. The endpoint is now
	// being executed the second time. So because we have our transaction response
	// tracked, the execution counts should still be one.
	{
		r, err := http.NewRequest("GET", "/test-path", nil)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
		r.Header.Add(TransactionIDHeader, "my-very-valid-test-transaction-id")
		w := httptest.NewRecorder()

		newServer.Config().Router.ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Fatal("expected", http.StatusOK, "got", w.Code)
		}

		decoderExecuted := e.(*testEndpoint).decoderExecuted
		if decoderExecuted != 1 {
			t.Fatal("expected", 1, "got", decoderExecuted)
		}
		endpointExecuted := e.(*testEndpoint).endpointExecuted
		if endpointExecuted != 1 {
			t.Fatal("expected", 1, "got", endpointExecuted)
		}
		encoderExecuted := e.(*testEndpoint).encoderExecuted
		if encoderExecuted != 1 {
			t.Fatal("expected", 1, "got", encoderExecuted)
		}

		if w.Body.String() != "test-response-1" {
			t.Fatal("expected", "test-response-1", "got", w.Body.String())
		}
	}
}

func Test_Transaction_InvalidIDGiven(t *testing.T) {
	e := testNewEndpoint(t)

	config := DefaultConfig()
	config.Endpoints = []Endpoint{e}
	config.ErrorEncoder = func(ctx context.Context, serverError error, w http.ResponseWriter) {
		w.WriteHeader(http.StatusInternalServerError)
	}
	newServer, err := New(config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	newServer.Boot()
	defer newServer.Shutdown()

	// Here we make a request against our test endpoint. The endpoint is provided
	// with an invalid transaction ID. The server's error encoder returns status
	// code 500 on all errors.
	{
		r, err := http.NewRequest("GET", "/test-path", nil)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
		r.Header.Add(TransactionIDHeader, "--my-invalid-transaction-id--")
		w := httptest.NewRecorder()

		newServer.Config().Router.ServeHTTP(w, r)

		if w.Code != http.StatusInternalServerError {
			t.Fatal("expected", http.StatusInternalServerError, "got", w.Code)
		}

		if !strings.Contains(w.Body.String(), "invalid transaction ID: does not match") {
			t.Fatal("expected", "invalid transaction ID: does not match", "got", w.Body.String())
		}
	}
}
