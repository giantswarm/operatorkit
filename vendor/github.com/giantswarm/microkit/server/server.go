// Package server provides a server implementation to connect network transport
// protocols and service business logic by defining server endpoints.
package server

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	kitendpoint "github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"github.com/tylerb/graceful"
	"golang.org/x/net/context"

	microerror "github.com/giantswarm/microkit/error"
	"github.com/giantswarm/microkit/logger"
	micrologger "github.com/giantswarm/microkit/logger"
	"github.com/giantswarm/microkit/transaction"
	microtransaction "github.com/giantswarm/microkit/transaction"
	transactionid "github.com/giantswarm/microkit/transaction/context/id"
	transactiontracked "github.com/giantswarm/microkit/transaction/context/tracked"
)

// Config represents the configuration used to create a new server object.
type Config struct {
	// Dependencies.

	// ErrorEncoder is the server's error encoder. This wraps the error encoder
	// configured by the client. Clients should not implement error logging in
	// here them self. This is done by the server itself. Clients must not
	// implement error response writing them self. This is done by the server
	// itself. Duplicated response writing will lead to runtime panics.
	ErrorEncoder kithttp.ErrorEncoder
	// Logger is the logger used to print log messages.
	Logger micrologger.Logger
	// Router is a HTTP handler for the server. The returned router will have all
	// endpoints registered that are listed in the endpoint collection.
	Router *mux.Router
	// TransactionResponder is the responder used to reply to requests using
	// persisted transaction results.
	TransactionResponder microtransaction.Responder

	// Settings.

	// Endpoints is the server's configured list of endpoints. These are the
	// custom endpoints configured by the client.
	Endpoints []Endpoint
	// HandlerWrapper is a wrapper provided to interact with the request on its
	// roots.
	HandlerWrapper func(h http.Handler) http.Handler
	// ListenAddress is the address the server is listening on.
	ListenAddress string
	// RequestFuncs is the server's configured list of request functions. These
	// are the custom request functions configured by the client.
	RequestFuncs []kithttp.RequestFunc
	// ServiceName is the name of the micro-service implementing the microkit
	// server. This is used for logging and instrumentation.
	ServiceName string
	// TLSCAFile is the file path to the certificate root CA file, if any.
	TLSCAFile string
	// TLSKeyFilePath is the file path to the certificate public key file, if any.
	TLSCrtFile string
	// TLSKeyFilePath is the file path to the certificate private key file, if
	// any.
	TLSKeyFile string
	// Viper is a configuration management object.
	Viper *viper.Viper
}

// DefaultConfig provides a default configuration to create a new server object
// by best effort.
func DefaultConfig() Config {
	var err error

	var loggerService micrologger.Logger
	{
		loggerConfig := micrologger.DefaultConfig()
		loggerService, err = micrologger.New(loggerConfig)
		if err != nil {
			panic(err)
		}
	}

	var responderService microtransaction.Responder
	{
		responderConfig := microtransaction.DefaultResponderConfig()
		responderService, err = microtransaction.NewResponder(responderConfig)
		if err != nil {
			panic(err)
		}
	}

	return Config{
		// Dependencies.
		ErrorEncoder:         func(ctx context.Context, serverError error, w http.ResponseWriter) {},
		Logger:               loggerService,
		Router:               mux.NewRouter(),
		TransactionResponder: responderService,

		// Settings.
		Endpoints:      nil,
		HandlerWrapper: func(h http.Handler) http.Handler { return h },
		ListenAddress:  "http://127.0.0.1:8000",
		RequestFuncs:   []kithttp.RequestFunc{},
		ServiceName:    "microkit",
		TLSCAFile:      "",
		TLSCrtFile:     "",
		TLSKeyFile:     "",
		Viper:          viper.New(),
	}
}

// New creates a new configured server object.
func New(config Config) (Server, error) {
	// Dependencies.
	if config.ErrorEncoder == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "error encoder must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "logger must not be empty")
	}
	if config.Router == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "router must not be empty")
	}
	if config.TransactionResponder == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "transaction responder must not be empty")
	}

	// Settings.
	if config.Endpoints == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "endpoints must not be empty")
	}
	if config.HandlerWrapper == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "handler wrapper must not be empty")
	}
	if config.ListenAddress == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "listen address must not be empty")
	}
	if config.RequestFuncs == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "request funcs must not be empty")
	}
	if config.ServiceName == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "service name must not be empty")
	}
	if config.TLSCrtFile == "" && config.TLSKeyFile != "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "TLS public key must not be empty")
	}
	if config.TLSCrtFile != "" && config.TLSKeyFile == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "TLS private key must not be empty")
	}

	listenURL, err := url.Parse(config.ListenAddress)
	if err != nil {
		return nil, microerror.MaskAnyf(invalidConfigError, err.Error())
	}

	newServer := &server{
		// Dependencies.
		errorEncoder:         config.ErrorEncoder,
		logger:               config.Logger,
		router:               config.Router,
		transactionResponder: config.TransactionResponder,

		// Internals.
		bootOnce:     sync.Once{},
		config:       config,
		httpServer:   nil,
		listenURL:    listenURL,
		shutdownOnce: sync.Once{},

		// Settings.
		endpoints:      config.Endpoints,
		handlerWrapper: config.HandlerWrapper,
		requestFuncs:   config.RequestFuncs,
		serviceName:    config.ServiceName,
		tlsCAFile:      config.TLSCAFile,
		tlsCrtFile:     config.TLSCrtFile,
		tlsKeyFile:     config.TLSKeyFile,
	}

	return newServer, nil
}

// server manages the transport logic and endpoint registration.
type server struct {
	// Dependencies.
	errorEncoder         kithttp.ErrorEncoder
	logger               logger.Logger
	router               *mux.Router
	transactionResponder transaction.Responder

	// Internals.
	bootOnce     sync.Once
	config       Config
	httpServer   *graceful.Server
	listenURL    *url.URL
	shutdownOnce sync.Once

	// Settings.
	endpoints      []Endpoint
	handlerWrapper func(h http.Handler) http.Handler
	requestFuncs   []kithttp.RequestFunc
	serviceName    string
	tlsCAFile      string
	tlsCrtFile     string
	tlsKeyFile     string
}

func (s *server) Boot() {
	s.bootOnce.Do(func() {
		s.router.NotFoundHandler = s.newNotFoundHandler()

		// Combine all options this server defines.
		options := []kithttp.ServerOption{
			kithttp.ServerBefore(s.requestFuncs...),
			kithttp.ServerErrorEncoder(s.newErrorEncoderWrapper()),
		}

		// We go through all endpoints this server defines and register them to the
		// router.
		for _, e := range s.endpoints {
			func(e Endpoint) {
				// Register all endpoints to the router depending on their HTTP methods and
				// request paths. The registered http.Handler is instrumented using
				// prometheus. We track counts of execution and duration it took to complete
				// the http.Handler.
				s.router.Methods(e.Method()).Path(e.Path()).Handler(s.handlerWrapper(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					ctx, err := s.newRequestContext(w, r)
					if err != nil {
						s.newErrorEncoderWrapper()(ctx, err, w)
						return
					}

					responseWriter, err := s.newResponseWriter(w)
					if err != nil {
						s.newErrorEncoderWrapper()(ctx, err, w)
						return
					}

					// Here we define the metrics labels. These will be used to instrument
					// the current request. This defered callback is initialized with the
					// timestamp of the beginning of the execution and will be executed at
					// the very end of the request. When it is executed we know all
					// necessary information to instrument the complete request, including
					// its response status code.
					defer func(t time.Time) {
						endpointCode := strconv.Itoa(responseWriter.StatusCode())
						endpointMethod := strings.ToLower(e.Method())
						endpointName := strings.Replace(e.Name(), "/", "_", -1)

						s.logger.Log("code", endpointCode, "endpoint", e.Name(), "method", endpointMethod, "path", r.URL.Path)

						endpointTotal.WithLabelValues(endpointCode, endpointMethod, endpointName).Inc()
						endpointTime.WithLabelValues(endpointCode, endpointMethod, endpointName).Set(float64(time.Since(t) / time.Millisecond))
					}(time.Now())

					// Wrapp the custom implementations of the endpoint specific business
					// logic.
					wrappedDecoder := s.newDecoderWrapper(e, responseWriter)
					wrappedEndpoint := s.newEndpointWrapper(e)
					wrappedEncoder := s.newEncoderWrapper(e, responseWriter)

					// Now we execute the actual go-kit endpoint handler.
					kithttp.NewServer(
						ctx,
						wrappedEndpoint,
						wrappedDecoder,
						wrappedEncoder,
						options...,
					).ServeHTTP(responseWriter, r)
				})))
			}(e)
		}

		// Register prometheus metrics endpoint.
		s.router.Path("/metrics").Handler(promhttp.Handler())

		// Register the router which has all of the configured custom endpoints
		// registered.
		s.httpServer = &graceful.Server{
			NoSignalHandling: true,
			Server: &http.Server{
				Addr:    s.listenURL.Host,
				Handler: s.router,
			},
			Timeout: 3 * time.Second,
		}

		go func() {
			if s.listenURL.Scheme == "https" {
				tlsConfig, err := s.newTLSConfig()
				if err != nil {
					panic(err)
				}
				s.logger.Log("debug", "running HTTPS server with TLS config")
				err = s.httpServer.ListenAndServeTLSConfig(tlsConfig)
				if err != nil {
					panic(err)
				}
			} else {
				s.logger.Log("debug", "running HTTP server")
				err := s.httpServer.ListenAndServe()
				if err != nil {
					panic(err)
				}
			}
		}()
	})
}

func (s *server) Config() Config {
	return s.config
}

func (s *server) Shutdown() {
	s.shutdownOnce.Do(func() {
		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			// Stop the HTTP server gracefully and wait some time for open connections
			// to be closed. Then force it to be stopped.
			s.httpServer.Stop(s.httpServer.Timeout)
			<-s.httpServer.StopChan()
			wg.Done()
		}()

		wg.Wait()
	})
}

func (s *server) newErrorEncoderWrapper() kithttp.ErrorEncoder {
	return func(ctx context.Context, serverError error, w http.ResponseWriter) {
		var err error

		// At first we have to set the content type of the actual error response. If
		// we would set it at the end we would set a trailing header that would not
		// be recognized by most of the clients out there. This is because in the
		// next call to the errorEncoder below the client's implementation of the
		// errorEncoder probably writes the status code header, which marks the
		// beginning of trailing headers in HTTP.
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		// Create the microkit specific response error, which acts as error wrapper
		// within the client's error encoder. It is used to propagate response codes
		// and messages, so we can use them below.
		var responseError ResponseError
		{
			responseConfig := DefaultResponseErrorConfig()
			responseConfig.Underlying = serverError
			responseError, err = NewResponseError(responseConfig)
			if err != nil {
				panic(err)
			}
		}

		rw, err := s.newResponseWriter(w)
		if err != nil {
			panic(err)
		}

		// Run the custom error encoder. This is used to let the implementing
		// microservice do something with errors occured during runtime. Things like
		// writing specific HTTP status codes to the given response writer or
		// writing data to the response body can be done.
		s.errorEncoder(ctx, responseError, rw)

		// Log the error and its errgo trace. This is really useful for debugging.
		errDomain := errorDomain(serverError)
		errMessage := errorMessage(serverError)
		errTrace := errorTrace(serverError)
		s.logger.Log("error", map[string]string{"domain": errDomain, "message": errMessage, "trace": errTrace})

		// Emit metrics about the occured errors. That way we can feed our
		// instrumentation stack to have nice dashboards to get a picture about the
		// general system health.
		errorTotal.WithLabelValues(errDomain).Inc()

		// Write the actual response body in case no response was already written
		// inside the error encoder.
		if !rw.HasWritten() {
			json.NewEncoder(rw).Encode(map[string]interface{}{
				"code":  responseError.Code(),
				"error": responseError.Message(),
				"from":  s.serviceName,
			})
		}
	}
}

// newDecoderWrapper creates a new wrappeed endpoint decoder. E.g. here we wrap
// the endpoint's decoder. We check if there is a transaction response being
// tracked. In this case we reply to the current request with the tracked
// information of the transaction response. After that the endpoint and encoder
// is not executed. Only response functions, if any, will be executed as usual.
// If there is no transaction response being tracked, the request is processed
// normally. This means that the usual execution of the endpoints decoder,
// endpoint and encoder takes place.
func (s *server) newDecoderWrapper(e Endpoint, responseWriter ResponseWriter) kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		tracked, ok := transactiontracked.FromContext(ctx)
		if !ok {
			return nil, microerror.MaskAnyf(invalidContextError, "tracked must not be empty")
		}
		if tracked {
			transactionID, ok := transactionid.FromContext(ctx)
			if !ok {
				return nil, microerror.MaskAnyf(invalidContextError, "transaction ID must not be empty")
			}
			err := s.transactionResponder.Reply(ctx, transactionID, responseWriter)
			if err != nil {
				return nil, microerror.MaskAny(err)
			}

			return nil, nil
		}

		request, err := e.Decoder()(ctx, r)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}

		return request, nil
	}
}

// newEncoderWrapper creates a new wrapped endpoint encoder. E.g. here we wrap
// the endpoint's decoder. In case the response of the current request is known
// to be tracked, we skip the execution of the actual endpoint. We rely on the
// wrapped decoder above, which already prepared the reply of the current
// request. If there is no transaction response being tracked, we execute the
// actual encoder as usual. Its response is being tracked in case a transaction
// ID is provided in the given request context. This tracked transaction
// response is used to reply to upcoming requests that provide the same
// transaction ID again.
func (s *server) newEncoderWrapper(e Endpoint, responseWriter ResponseWriter) kithttp.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		tracked, ok := transactiontracked.FromContext(ctx)
		if !ok {
			return microerror.MaskAnyf(invalidContextError, "tracked must not be empty")
		}
		if tracked {
			return nil
		}

		err := e.Encoder()(ctx, w, response)
		if err != nil {
			return microerror.MaskAny(err)
		}

		transactionID, ok := transactionid.FromContext(ctx)
		if !ok {
			// In case the response is not already tracked, but there is no
			// transaction ID, we cannot track it at all. So we return here.
			return nil
		}
		err = s.transactionResponder.Track(ctx, transactionID, responseWriter)
		if err != nil {
			return microerror.MaskAny(err)
		}

		return nil
	}
}

// newEndpointWrapper creates a new wrapped endpoint function. E.g. here we wrap
// the actual endpoint, the business logic. In case the response of the current
// request is known to be tracked, we skip the execution of the actual endpoint.
// We rely on the wrapped decoder above, which already prepared the reply of the
// current request. If there is no transaction response being tracked, we
// execute the actual endpoint as usual.
func (s *server) newEndpointWrapper(e Endpoint) kitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tracked, ok := transactiontracked.FromContext(ctx)
		if !ok {
			return nil, microerror.MaskAnyf(invalidContextError, "tracked must not be empty")
		}
		if tracked {
			return nil, nil
		}

		// Prepare the actual endpoint depending on the provided middlewares of the
		// endpoint implementation. There might be cases in which there are none or
		// only one middleware. The go-kit interface is not that nice so we need to
		// make it fit here.
		endpoint := e.Endpoint()
		middlewares := e.Middlewares()
		if len(middlewares) == 1 {
			endpoint = kitendpoint.Chain(middlewares[0])(endpoint)
		}
		if len(middlewares) > 1 {
			endpoint = kitendpoint.Chain(middlewares[0], middlewares[1:]...)(endpoint)
		}
		response, err := endpoint(ctx, request)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}

		return response, nil
	}
}

// newNotFoundHandler returns an HTTP handler that represents our custom not
// found handler. Here we take care about logging, metrics and a proper
// response.
func (s *server) newNotFoundHandler() http.Handler {
	return http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the error and its message. This is really useful for debugging.
		errDomain := errorDomain(nil)
		errMessage := fmt.Sprintf("not found: %s %s", r.Method, r.URL.Path)
		errTrace := ""
		s.logger.Log("error", map[string]string{"domain": errDomain, "message": errMessage, "trace": errTrace})

		// This defered callback will be executed at the very end of the request.
		defer func(t time.Time) {
			endpointCode := strconv.Itoa(http.StatusNotFound)
			endpointMethod := strings.ToLower(r.Method)
			endpointName := "notfound"

			endpointTotal.WithLabelValues(endpointCode, endpointMethod, endpointName).Inc()
			endpointTime.WithLabelValues(endpointCode, endpointMethod, endpointName).Set(float64(time.Since(t) / time.Millisecond))

			errorTotal.WithLabelValues(errDomain).Inc()
		}(time.Now())

		// Write the actual response body.
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":  CodeResourceNotFound,
			"error": errMessage,
			"from":  s.serviceName,
		})
	}))
}

// newRequestContext creates a new request context and enriches it with request
// relevant information. E.g. here we put the HTTP X-Transaction-ID header into
// the request context, if any. We also check if there is a transaction response
// already tracked for the given transaction ID. This information is then stored
// within the given request context as well. Note that we initialize the
// information about the tracked state of the transaction response with false,
// to always have a valid state available within the request context.
func (s *server) newRequestContext(w http.ResponseWriter, r *http.Request) (context.Context, error) {
	ctx := context.Background()
	ctx = transactiontracked.NewContext(ctx, false)

	transactionID := r.Header.Get(TransactionIDHeader)
	if transactionID == "" {
		return ctx, nil
	}

	if !IsValidTransactionID(transactionID) {
		return nil, microerror.MaskAnyf(invalidTransactionIDError, "does not match %s", TransactionIDRegEx.String())
	}

	ctx = transactionid.NewContext(ctx, transactionID)

	exists, err := s.transactionResponder.Exists(ctx, transactionID)
	if err != nil {
		return nil, microerror.MaskAny(err)
	}
	ctx = transactiontracked.NewContext(ctx, exists)

	return ctx, nil
}

// newResponseWriter creates a new wrapped HTTP response writer. E.g. here we
// create a new wrapper for the http.ResponseWriter of the current request. We
// inject it into the called http.Handler so it can track the status code we are
// interested in. It will help us gathering the response status code after it
// was written by the underlying http.ResponseWriter.
func (s *server) newResponseWriter(w http.ResponseWriter) (ResponseWriter, error) {
	responseConfig := DefaultResponseWriterConfig()
	responseConfig.ResponseWriter = w
	responseWriter, err := NewResponseWriter(responseConfig)
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	return responseWriter, nil
}

func (s *server) newTLSConfig() (*tls.Config, error) {
	var err error

	var roots *x509.CertPool
	if s.tlsCAFile != "" {
		roots, err = x509.SystemCertPool()
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
		b, err := ioutil.ReadFile(s.tlsCAFile)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
		ok := roots.AppendCertsFromPEM(b)
		if !ok {
			return nil, microerror.MaskAny(fmt.Errorf("could not load root CA: '%s'", s.tlsCAFile))
		}
		s.logger.Log("debug", fmt.Sprintf("found TLS root CA file '%s'", s.tlsCAFile))
	}

	var certs []tls.Certificate
	if s.tlsCrtFile != "" && s.tlsKeyFile != "" {
		c, err := tls.LoadX509KeyPair(s.tlsCrtFile, s.tlsKeyFile)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
		s.logger.Log("debug", fmt.Sprintf("found TLS public key file '%s' and private key file '%s'", s.tlsCrtFile, s.tlsKeyFile))
		certs = append(certs, c)
	}

	tlsConfig := &tls.Config{
		Certificates: certs,
		RootCAs:      roots,
	}

	return tlsConfig, nil
}
