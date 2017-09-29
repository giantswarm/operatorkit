package informer

import (
	"context"
	"fmt"
	"time"

	"github.com/cenk/backoff"
	"github.com/giantswarm/microerror"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
)

const (
	// ResyncPeriod is the interval at which the Informer cache is invalidated,
	// and the lister function is called.
	ResyncPeriod = 1 * time.Minute
)

// Config represents the configuration used to create a new Informer.
type Config struct {
	// Dependencies.
	BackOff    backoff.BackOff
	RestClient rest.Interface

	// Settings.
	ResyncPeriod time.Duration
}

// DefaultConfig provides a default configuration to create a new by best
// effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		BackOff:    nil,
		RestClient: nil,

		// Settings.
		ResyncPeriod: ResyncPeriod,
	}
}

type Informer struct {
	// Dependencies.
	backOff    backoff.BackOff
	restClient rest.Interface

	// Settings.
	resyncPeriod time.Duration
}

// New creates a new Informer.
func New(config Config) (*Informer, error) {
	// Dependencies.
	if config.BackOff == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.BackOff must not be empty")
	}
	if config.RestClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.RestClient must not be empty")
	}

	// Settings.
	if config.ResyncPeriod == 0 {
		return nil, microerror.Maskf(invalidConfigError, "config.ResyncPeriod must not be empty")
	}

	newInformer := &Informer{
		// Settings.
		backOff:    config.BackOff,
		restClient: config.RestClient,

		// Settings.
		resyncPeriod: config.ResyncPeriod,
	}

	return newInformer, nil
}

// Watch only watches objects using a stream decoder. Afer the resync period the
// active watch is closed and a new stream decoder watches the API again. This
// mechanism has potential to not recognize delete events of objects that do not
// use finalizers. This should be safe though when using Watch in combination
// with the operatorkit reconciler framework since it uses finalizers for all
// event objects.
func (i *Informer) Watch(ctx context.Context, p WatchEndpointProvider, f ZeroObjectFactory) (chan watch.Event, chan watch.Event, chan error) {
	done := make(chan struct{}, 1)

	deleteChan := make(chan watch.Event, 1)
	updateChan := make(chan watch.Event, 1)
	errChan := make(chan error, 1)

	// Having a separate method for creating the watcher and using it to fetch
	// events is more convenient when controling the loop and select flows because
	// we can just return and do not need to use ugly labels.
	fetchEvents := func(canceler chan struct{}) {
		stream, err := i.restClient.Get().AbsPath(p.WatchEndpoint()).Stream()
		if err != nil {
			errChan <- microerror.Mask(err)
			return
		}
		watcher := watch.NewStreamWatcher(newDecoder(stream, f))

		defer watcher.Stop()

		for {
			select {
			case <-done:
				return
			case <-canceler:
				return
			case event, ok := <-watcher.ResultChan():
				if ok {
					switch event.Type {
					case watch.Added, watch.Modified:
						updateChan <- event
					case watch.Deleted:
						deleteChan <- event
					case watch.Error:
						errChan <- microerror.Maskf(invalidEventError, "%#v", event)
					default:
						errChan <- microerror.Maskf(invalidEventError, "%#v", event)
					}
				} else {
					fmt.Printf("no event found\n")
					return
				}
			}
		}
	}

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				canceler := make(chan struct{})
				go fetchEvents(canceler)
				<-time.After(i.resyncPeriod)
				close(canceler)
			}
		}
	}()

	go func() {
		<-ctx.Done()

		close(done)

		close(deleteChan)
		close(updateChan)
		close(errChan)
	}()

	return deleteChan, updateChan, errChan
}
