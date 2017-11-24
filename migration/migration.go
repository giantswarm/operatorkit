package migration

import (
	"context"
	"time"

	"github.com/cenk/backoff"
	"github.com/giantswarm/microerror"
)

// Process takes a list of ordered migrators and runs them sequentially. The
// order of migrators should never change once used to guarantee a constantly
// stable migration path.
func Process(ctx context.Context, migrators []Migrator) error {
	for _, m := range migrators {
		o := func() error {
			err := m.Init()
			if err != nil {
				return microerror.Mask(err)
			}

			eventChan, err := m.List(ctx)
			if err != nil {
				return microerror.Mask(err)
			}

			for {
				select {
				case <-time.After(time.Second):
					return nil
				case e := <-eventChan:
					t, err := m.Transform(e.Object)
					if err != nil {
						return microerror.Mask(err)
					}

					err = m.Create(t)
					if err != nil {
						return microerror.Mask(err)
					}

					err = m.Delete(e.Object)
					if err != nil {
						return microerror.Mask(err)
					}
				}
			}

			return nil
		}

		b := backoff.WithMaxTries(backoff.NewExponentialBackOff(), 3)
		err := backoff.Retry(o, b)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}
