// +build k8srequired

package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func Test_Informer_Collector(t *testing.T) {
	fmt.Println("setting up")
	if err := mustSetup(); err != nil {
		t.Fatal("expected", nil, "got", err)
	}
	defer mustTeardown()

	idOne := "al7qy"

	informer, err := newOperatorkitInformer(time.Second*2, time.Second*10)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	if err := createConfigMap(idOne); err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	metrics := []prometheus.Metric{}
	metricChan := make(chan prometheus.Metric)

	go func() {
		informer.Collect(metricChan)
	}()

	fmt.Println("deleting config map")
	if err := deleteConfigMap(idOne); err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	for metric := range metricChan {
		metrics = append(metrics, metric)
	}

	if len(metrics) != 1 {
		t.Fatal("expected", 1, "got", len(metrics))
	}
}
