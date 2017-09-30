package informer

import (
	"context"
	"sync"
	"time"

	"github.com/cenk/backoff"
	"github.com/giantswarm/microerror"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
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

	// Internals.
	cache *sync.Map

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

		// Internals.
		cache: &sync.Map{},

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
func (i *Informer) Watch(ctx context.Context, endpoint string, factory ZeroObjectFactory) (chan watch.Event, chan watch.Event, chan watch.Event, chan error) {
	done := make(chan struct{}, 1)
	eventChan := make(chan watch.Event, 1)

	createChan := make(chan watch.Event, 1)
	deleteChan := make(chan watch.Event, 1)
	updateChan := make(chan watch.Event, 1)
	errChan := make(chan error, 1)

	go func() {
		for {
			select {
			case <-done:
			case event := <-eventChan:
				switch event.Type {
				case watch.Added:
					err := i.cacheOrReleaseEvent(event, createChan)
					if err != nil {
						errChan <- microerror.Mask(err)
					}
				case watch.Deleted:
					err := i.uncacheAndReleaseEvent(event, deleteChan)
					if err != nil {
						errChan <- microerror.Mask(err)
					}
				case watch.Modified:
					updateChan <- event
				default:
					errChan <- microerror.Maskf(invalidEventError, "%#v", event)
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				ctx, cancelFunc := context.WithCancel(ctx)
				go func() {
					for {
						select {
						case <-ctx.Done():
							return
						default:
							err := i.streamEvents(ctx, endpoint, factory, eventChan)
							if err != nil {
								errChan <- microerror.Mask(err)
							}
						}
					}
				}()

				<-time.After(i.resyncPeriod)

				i.releaseCachedEvents(ctx, updateChan)
				cancelFunc()
			}
		}
	}()

	go func() {
		<-ctx.Done()

		close(done)
		close(eventChan)

		close(createChan)
		close(deleteChan)
		close(updateChan)
		close(errChan)
	}()

	return createChan, deleteChan, updateChan, errChan
}

func (i *Informer) cacheOrReleaseEvent(event watch.Event, createChan chan watch.Event) error {
	k, err := cache.MetaNamespaceKeyFunc(event.Object)
	if err != nil {
		return microerror.Mask(err)
	}

	_, ok := i.cache.Load(k)
	if !ok {
		createChan <- event
	}

	i.cache.Store(k, event)

	return nil
}

func (i *Informer) uncacheAndReleaseEvent(event watch.Event, deleteChan chan watch.Event) error {
	deleteChan <- event

	k, err := cache.MetaNamespaceKeyFunc(event.Object)
	if err != nil {
		return microerror.Mask(err)
	}

	i.cache.Delete(k)

	return nil
}

func (i *Informer) releaseCachedEvents(ctx context.Context, updateChan chan watch.Event) {
	i.cache.Range(func(k, v interface{}) bool {
		select {
		case <-ctx.Done():
			return false
		default:
			updateChan <- v.(watch.Event)
		}
		return true
	})
}

func (i *Informer) streamEvents(ctx context.Context, endpoint string, factory ZeroObjectFactory, eventChan chan watch.Event) error {
	stream, err := i.restClient.Get().AbsPath(endpoint).Stream()
	if err != nil {
		return microerror.Mask(err)
	}
	watcher := watch.NewStreamWatcher(newDecoder(stream, factory))

	defer watcher.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case event, ok := <-watcher.ResultChan():
			if ok {
				eventChan <- event
			} else {
				return nil
			}
		}
	}

	return nil
}
