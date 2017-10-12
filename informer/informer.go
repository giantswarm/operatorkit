// package informer provides primitives to watch event objects from the
// Kubernetes API in a deterministic way. The following conditions are
// guaranteed by the watcher.
//
//     - The informer is able to watch all kinds of objects as soon as a proper
//       watch endpoint and a factory implementing ZeroObjectFactory is given.
//     - Events for objects that are created, deleted or updated are dispatched
//       immediately.
//     - Events for objects that are created or updated are dispatched via the
//       same channel. The informer cannot distinguish between a created or
//       updated event object.
//     - Events for objects that are alredy cached are not dispatched during the
//       configured resync period, but periodically after it.
//     - Events for objects are never dispatched twice.
//     - Events for objects can be dispatched in a rate limitted way, if
//       configured accordingly.
//
package informer

import (
	"context"
	"sync"
	"time"

	"github.com/cenk/backoff"
	"github.com/giantswarm/microerror"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

const (
	// DefaultRateWait is the default value for the RateWait setting. See Config
	// for more information.
	DefaultRateWait = 0 * time.Second
	// DefaultResyncPeriod is the default value for the ResyncPeriod setting. See
	// Config for more information.
	DefaultResyncPeriod = 1 * time.Minute
)

// Config represents the configuration used to create a new Informer.
type Config struct {
	// Dependencies.

	BackOff        backoff.BackOff
	WatcherFactory WatcherFactory

	// Settings.

	// RateWait provides configuration for some kind of rate limitting. The
	// informer watch provides events via the update channel every ResyncPeriod.
	// This triggers the release of update events. RateWait is the time to wait
	// between released events.
	RateWait time.Duration
	// ResyncPeriod is the time to wait before releasing update events again.
	ResyncPeriod time.Duration
}

// DefaultConfig provides a default configuration to create a new by best
// effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		BackOff:        nil,
		WatcherFactory: nil,

		// Settings.
		RateWait:     DefaultRateWait,
		ResyncPeriod: DefaultResyncPeriod,
	}
}

// Informer implements primitives to watch event objects from the Kubernetes API
// in a deterministic way.
type Informer struct {
	// Dependencies.
	backOff        backoff.BackOff
	watcherFactory WatcherFactory

	// Internals.
	cache       *sync.Map
	initializer chan struct{}

	// Settings.
	rateWait     time.Duration
	resyncPeriod time.Duration
}

// New creates a new Informer.
func New(config Config) (*Informer, error) {
	// Dependencies.
	if config.BackOff == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.BackOff must not be empty")
	}
	if config.WatcherFactory == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.WatcherFactory must not be empty")
	}

	// Settings.
	if config.ResyncPeriod == 0 {
		return nil, microerror.Maskf(invalidConfigError, "config.ResyncPeriod must not be empty")
	}

	newInformer := &Informer{
		// Settings.
		backOff:        config.BackOff,
		watcherFactory: config.WatcherFactory,

		// Internals.
		cache:       &sync.Map{},
		initializer: make(chan struct{}),

		// Settings.
		rateWait:     config.RateWait,
		resyncPeriod: config.ResyncPeriod,
	}

	return newInformer, nil
}

// Watch only watches objects using a stream decoder. Afer the resync period the
// active watch is closed and a new stream decoder watches the API again. This
// mechanism has a very small potential to not recognize delete events of
// objects that do not use finalizers.
//
// Watch takes a context as first argument which can be used to cancel the
// watch. The second argument provided is the raw API endpoint path the watcher
// is using to fetch event objects. The third argument is a custom object
// factory to create zero objects on demand. The watcher is using this to decode
// event objects.
//
// Watch returns channels for delete, update and error events, in this order.
// Events will be dispatched as soon as they happen.
//
// That the resync period configured for the informer will trigger periodic
// updates of event objects via the update channel.
func (i *Informer) Watch(ctx context.Context) (chan watch.Event, chan watch.Event, chan error) {
	done := make(chan struct{}, 1)
	eventChan := make(chan watch.Event, 1)

	deleteChan := make(chan watch.Event, 1)
	updateChan := make(chan watch.Event, 1)
	errChan := make(chan error, 1)

	go func() {
		for {
			select {
			case <-done:
				return
			case event, ok := <-eventChan:
				if !ok {
					return
				}

				switch event.Type {
				case watch.Added:
					err := i.cacheAndSendIfNotExists(event, updateChan)
					if err != nil {
						errChan <- microerror.Mask(err)
					}
				case watch.Deleted:
					err := i.uncacheAndSend(event, deleteChan)
					if err != nil {
						errChan <- microerror.Mask(err)
					}
				case watch.Modified:
					err := i.cacheAndSend(event, deleteChan, updateChan)
					if err != nil {
						errChan <- microerror.Mask(err)
					}
				default:
					errChan <- microerror.Maskf(invalidEventError, "%#v", event)
				}
			}
		}
	}()

	go func() {
		// Here we fill the informer cache initially and release the event objects
		// the very first time after the program started. This is a special case and
		// guarantees any configured rate limitting is properly done.
		{
			err := i.fillCache(ctx, eventChan)
			if err != nil {
				errChan <- microerror.Mask(err)
			}
			close(i.initializer)
			i.sendCachedEvents(ctx, deleteChan, updateChan, errChan)
		}

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
							err := i.streamEvents(ctx, eventChan)
							if err != nil {
								errChan <- microerror.Mask(err)
							}
						}
					}
				}()

				<-time.After(i.resyncPeriod)

				i.sendCachedEvents(ctx, deleteChan, updateChan, errChan)
				cancelFunc()
			}
		}
	}()

	go func() {
		<-ctx.Done()

		close(done)
		close(eventChan)

		close(deleteChan)
		close(updateChan)
		close(errChan)
	}()

	return deleteChan, updateChan, errChan
}

// cacheAndSend stores the provided event object in the informer cache and
// dispatches it based on its properties. cacheAndSend sends the provided event
// object to the provided update channel in case the event object has no
// deletion timestamp. In case the deletion timestamp of the provided event
// object is non-nil, it is send to the provided delete channel.
func (i *Informer) cacheAndSend(event watch.Event, deleteChan, updateChan chan watch.Event) error {
	k, err := cache.MetaNamespaceKeyFunc(event.Object)
	if err != nil {
		return microerror.Mask(err)
	}
	i.cache.Store(k, event)

	m, err := meta.Accessor(event.Object)
	if err != nil {
		return microerror.Mask(err)
	}
	t := m.GetDeletionTimestamp()
	if t == nil {
		updateChan <- event
	} else {
		deleteChan <- event
	}

	return nil
}

// cacheAndSendIfNotExists handles watch.ADDED events. These events can happen
// because of different reasons.
//
//     - The watcher may receives a new event object because a new object was
//       created in the API.
//     - The watcher may syncs the very first time on programm start.
//     - The watcher may resyncs and receives an event object we already know.
//
// In case the provided event object does not exist in the informer cache, this
// means we send it to the provided update channel because it should be
// reconciled. The reconciliation has to make sure the event object is created
// and/or updated accordingly. In any case cacheAndSendIfNotExists adds the
// provided event object to the informer cache.
func (i *Informer) cacheAndSendIfNotExists(event watch.Event, updateChan chan watch.Event) error {
	k, err := cache.MetaNamespaceKeyFunc(event.Object)
	if err != nil {
		return microerror.Mask(err)
	}

	_, ok := i.cache.Load(k)
	if !ok && i.isCachedFilled() {
		updateChan <- event
	}

	i.cache.Store(k, event)

	return nil
}

// fillCache is similar to streamEvents but is only used during informer cache
// initialization. As soon as the watcher does not receive any event objects
// anymore, the cache is filled and the usual event watching process can begin.
func (i *Informer) fillCache(ctx context.Context, eventChan chan watch.Event) error {
	watcher, err := i.watcherFactory()
	if err != nil {
		return microerror.Mask(err)
	}

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
		case <-time.After(time.Second):
			return nil
		}
	}

	return nil
}

// isCachedFilled checks whether the informer cache is filled.
func (i *Informer) isCachedFilled() bool {
	select {
	case <-i.initializer:
		return true
	default:
		// fall thorugh
	}

	return false
}

// sendCachedEvents sends all cached event objects to the provided delete or
// update channel. sendCachedEvents sends the provided event object to the
// provided update channel in case the event object has no deletion timestamp.
// In case the deletion timestamp of the provided event object is non-nil, it is
// send to the provided delete channel. The release process may be rate limitted
// by the rate wait configuration of the informer. Then the release sleeps for
// the configured duration before releasing the next event object.
func (i *Informer) sendCachedEvents(ctx context.Context, deleteChan, updateChan chan watch.Event, errChan chan error) {
	// useRateWait is used to not apply the configured rate wait on the very first
	// event object. This is done to not wait any additional time before releasing
	// the first event object after the configured resync period.
	var useRateWait bool

	i.cache.Range(func(k, v interface{}) bool {
		e := v.(watch.Event)

		if useRateWait && i.rateWait != 0 {
			<-time.After(i.rateWait)
		}
		useRateWait = true

		select {
		case <-ctx.Done():
			return false
		default:
			m, err := meta.Accessor(e.Object)
			if err != nil {
				errChan <- microerror.Mask(err)
			} else {
				t := m.GetDeletionTimestamp()
				if t == nil {
					updateChan <- e
				} else {
					deleteChan <- e
				}
			}
		}

		return true
	})
}

// streamEvents creates a new watcher and sends event objects the watcher
// receives. It may happen that the watcher gets closed automatically, e.g. due
// to connection issues. As soon as the watcher gets closed or the watch gets
// canceled via the done channel of the provided context, streamEvents returns
// and stops blocking.
func (i *Informer) streamEvents(ctx context.Context, eventChan chan watch.Event) error {
	watcher, err := i.watcherFactory()
	if err != nil {
		return microerror.Mask(err)
	}

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

// uncacheAndSend sends the received event to the provided delete channel and
// removes the event object from the internal informer cache.
func (i *Informer) uncacheAndSend(event watch.Event, deleteChan chan watch.Event) error {
	deleteChan <- event

	k, err := cache.MetaNamespaceKeyFunc(event.Object)
	if err != nil {
		return microerror.Mask(err)
	}

	i.cache.Delete(k)

	return nil
}
