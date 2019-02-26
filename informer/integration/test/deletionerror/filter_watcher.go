// +build k8srequired

package deletionerror

import (
	"github.com/giantswarm/microerror"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/giantswarm/operatorkit/informer"
)

type FilterWatcherConfig struct {
	Watcher informer.Watcher
}

// FilterWatcher is a wrapping implementation of infomrmer.Watch used to filter
// watched events on demand. Filtering can be switched using SetDispatchEvents.
type FilterWatcher struct {
	watcher informer.Watcher

	dispatchEvents bool
}

func NewFilterWatcher(config FilterWatcherConfig) (*FilterWatcher, error) {
	if config.Watcher == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Watcher must not be empty", config)
	}

	w := &FilterWatcher{
		watcher: config.Watcher,

		dispatchEvents: true,
	}

	return w, nil
}

func (w *FilterWatcher) SetDispatchEvents(dispatchEvents bool) {
	w.dispatchEvents = dispatchEvents
}

func (w *FilterWatcher) Watch(listOptions metav1.ListOptions) (watch.Interface, error) {
	watchInterface, err := w.watcher.Watch(listOptions)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return watch.Filter(watchInterface, w.filterFunc), nil
}

func (w *FilterWatcher) filterFunc(e watch.Event) (watch.Event, bool) {
	return e, w.dispatchEvents
}
