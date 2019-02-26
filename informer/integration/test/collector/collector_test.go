// +build k8srequired

package collector

import (
	"context"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/giantswarm/operatorkit/informer/collector"
)

func Test_Informer_Collector_MultipleEvents(t *testing.T) {
	mustSetup()
	defer mustTeardown()

	var err error

	// We run the creation, update and deletion of a configmap multiple times to
	// simulate multiple occuring events in the informer watcher. When the
	// collector implementation using the informer watcher is correct and prevents
	// duplicates, no error should be produced gathering metrics below.
	{
		ctx, cancelFunc := context.WithCancel(context.Background())

		go func() {
			time.Sleep(5 * time.Second)
			cancelFunc()
		}()

		go func() {
			idOne := "al7qy"

			for {
				select {
				case <-ctx.Done():
					return
				default:
					err := createConfigMap(idOne)
					if err != nil {
						t.Fatal("expected", nil, "got", err)
					}
					err = updateConfigMap(idOne)
					if err != nil {
						t.Fatal("expected", nil, "got", err)
					}
					err = deleteConfigMap(idOne)
					if err != nil {
						t.Fatal("expected", nil, "got", err)
					}
				}
			}
		}()
	}

	// We create a Prometheus registry that registers and gathers metrics from the
	// collector implementation we want to ensure works properly.
	var r *prometheus.Registry
	{
		r = prometheus.NewRegistry()
	}

	// The collector implementation initialized below is used by the informer.
	var c *collector.Set
	{
		c, err = newInformerCollector()
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	// We register the collector in our test registry so the registry can gather
	// metrics from the collector. When the collector implementation produces e.g.
	// duplicated metrics in terms of duplicated label pairs, the Gather() call
	// below will return an error. Note the 26 errors reported in the example
	// below are caused during the 5 second period of creating, updating, and
	// deleting the same configmap over and over again during the test before
	// the fix was introduced.
	//
	//     --- FAIL: Test_Informer_Collector_MultipleEvents (6.62s)
	//         collector_test.go:93: expected <nil> got 26 error(s) occurred:
	//             * collected metric "operatorkit_informer_creation_timestamp" { label:<name:"kind" value:"" > label:<name:"name" value:"al7qy" > label:<name:"namespace" value:"test-informer-integration-collector" > gauge:<value:1.541631308e+09 > } was collected before with the same name and label values
	//             * collected metric "operatorkit_informer_creation_timestamp" { label:<name:"kind" value:"" > label:<name:"name" value:"al7qy" > label:<name:"namespace" value:"test-informer-integration-collector" > gauge:<value:1.541631308e+09 > } was collected before with the same name and label values
	//             ...
	//
	{
		err := r.Register(c)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
		_, err = r.Gather()
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}
}
