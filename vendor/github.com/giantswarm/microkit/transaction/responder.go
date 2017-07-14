package transaction

import (
	"context"
	"encoding/json"
	"fmt"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"
	microstorage "github.com/giantswarm/microkit/storage"
)

// ResponderConfig represents the configuration used to create a responder.
type ResponderConfig struct {
	// Dependencies.
	Logger  micrologger.Logger
	Storage microstorage.Service
}

// DefaultResponderConfig provides a default configuration to create a new
// responder by best effort.
func DefaultResponderConfig() ResponderConfig {
	var err error

	var loggerService micrologger.Logger
	{
		loggerConfig := micrologger.DefaultConfig()
		loggerService, err = micrologger.New(loggerConfig)
		if err != nil {
			panic(err)
		}
	}

	var storageService microstorage.Service
	{
		storageConfig := microstorage.DefaultConfig()
		storageService, err = microstorage.New(storageConfig)
		if err != nil {
			panic(err)
		}
	}

	config := ResponderConfig{
		// Dependencies.
		Logger:  loggerService,
		Storage: storageService,
	}

	return config
}

// NewResponder creates a new configured responder.
func NewResponder(config ResponderConfig) (Responder, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "logger must not be empty")
	}
	if config.Storage == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "storage must not be empty")
	}

	newResponder := &responder{
		// Dependencies.
		logger:  config.Logger,
		storage: config.Storage,
	}

	return newResponder, nil
}

type responder struct {
	// Dependencies.
	logger  micrologger.Logger
	storage microstorage.Service
}

func (r *responder) Exists(ctx context.Context, transactionID string) (bool, error) {
	key := responseKey("transaction", transactionID)
	exists, err := r.storage.Exists(ctx, key)
	if err != nil {
		return false, microerror.MaskAny(err)
	}

	return exists, nil
}

func (r *responder) Reply(ctx context.Context, transactionID string, rr ResponseReplier) error {
	// We search for the transaction response associated with the given
	// transaction ID. In case the client wanted to reply to a request, but we
	// cannot find the desired transaction response, we return an error.
	var response Response
	{
		key := responseKey("transaction", transactionID)
		res, err := r.storage.Search(ctx, key)
		if microstorage.IsNotFound(err) {
			return microerror.MaskAny(notFoundError)
		} else if err != nil {
			return microerror.MaskAny(err)
		}

		err = json.Unmarshal([]byte(res), &response)
		if err != nil {
			return microerror.MaskAny(err)
		}
	}

	// We found an existing transaction response. Here we write the actual HTTP
	// response using the tracked information of the transaction response.
	// headers.
	{
		for key, val := range response.Header {
			for _, h := range val {
				rr.Header().Add(key, h)
			}
		}

		rr.WriteHeader(response.Code)

		_, err := rr.Write([]byte(response.Body))
		if err != nil {
			return microerror.MaskAny(err)
		}

		key := responseKey("transaction", transactionID)
		r.logger.Log("debug", fmt.Sprintf("replied using transaction response with key '%s' and value '%#v'", key, response))
	}

	return nil
}

func (r *responder) Track(ctx context.Context, transactionID string, rt ResponseTracker) error {
	// At first we check if there does a transaction response already exist. We
	// only want to track responses once. So we return an error in case the client
	// would go to do duplicated work.
	{
		exists, err := r.Exists(ctx, transactionID)
		if err != nil {
			return microerror.MaskAny(err)
		}
		if exists {
			return microerror.MaskAnyf(alreadyExistsError, "transaction response for ID '%s' already exists", transactionID)
		}
	}

	// Here we use the given response tracker to create the transaction response
	// we actually want to persist. The string value created here is used below to
	// associate it with the given transaction ID within the underlying storage.
	var val string
	{
		response := Response{
			Body:   rt.BodyBuffer().String(),
			Code:   rt.StatusCode(),
			Header: rt.Header(),
		}
		b, err := json.Marshal(response)
		if err != nil {
			return microerror.MaskAny(err)
		}
		val = string(b)
	}

	// Now we have all information in place to store the transaction response.
	{
		key := responseKey("transaction", transactionID)
		err := r.storage.Create(ctx, key, val)
		if err != nil {
			return microerror.MaskAny(err)
		}

		r.logger.Log("debug", fmt.Sprintf("created transaction response with key '%s' and value '%s'", key, val))
	}

	return nil
}
