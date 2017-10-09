package informer

import (
	"github.com/giantswarm/microerror"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
)

func NewWatcherFactory(restClient rest.Interface, endpoint string, factory ZeroObjectFactory) WatcherFactory {
	return func() (watch.Interface, error) {
		stream, err := restClient.Get().AbsPath(endpoint).Stream()
		if err != nil {
			return nil, microerror.Mask(err)
		}

		watcher := watch.NewStreamWatcher(newDecoder(stream, factory))

		return watcher, nil
	}
}
