// +build k8srequired

package collector

import (
	"fmt"
	"testing"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/informer"
	"github.com/prometheus/client_golang/prometheus"
	prommodel "github.com/prometheus/client_model/go"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	grace int64 = 500
)

var (
	configMapName = "test-configmap"
	configmap     = &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: namespace,
		},
		Data: map[string]string{},
	}

	kindLabel      = "kind"
	nameLabel      = "name"
	namespaceLabel = "namespace"
	emptyValue     = ""
)

// Test_Informer_Collector_Basic is an integration test for basic collector operations.
// The test verifies the metrics emitted by the collector.
func Test_Informer_Collector_Basic(t *testing.T) {
	testCases := []struct {
		name string

		// Events to run before starting the collector.
		setupEvents []func(kubernetes.Interface) error
		// Events to run just after starting the collector.
		events []func(kubernetes.Interface) error

		// Descriptions of the metrics we expect to be emitted.
		expectedMetrics []*prometheus.Desc
	}{
		{
			name: "Test one create and one update event emits only one creation_timestamp metric",

			setupEvents: []func(kubernetes.Interface) error{
				// Create the configmap to force a create event.
				func(k8sClient kubernetes.Interface) error {
					if _, err := k8sClient.CoreV1().ConfigMaps(namespace).Create(configmap); err != nil {
						return microerror.Mask(err)
					}

					return nil
				},
			},

			events: []func(kubernetes.Interface) error{
				// Update the configmap to force an update event,
				// within the 1 second of the collector watching for events.
				func(k8sClient kubernetes.Interface) error {
					configmap.Data["foo"] = "bar"
					if _, err := k8sClient.CoreV1().ConfigMaps(namespace).Update(configmap); err != nil {
						return microerror.Mask(err)
					}
					delete(configmap.Data, "foo") // Be a good house guest.

					return nil
				},
			},

			// Check that only one metric is returned, as we only have one resource,
			// despite there being both a create and an update event.
			expectedMetrics: []*prometheus.Desc{
				informer.CreationTimestampDescription,
			},
		},

		{
			name: "Test X",

			setupEvents: []func(kubernetes.Interface) error{},

			events: []func(kubernetes.Interface) error{
				func(k8sClient kubernetes.Interface) error {
					if _, err := k8sClient.CoreV1().ConfigMaps(namespace).Create(configmap); err != nil {
						return microerror.Mask(err)
					}
					return nil
				},
				func(k8sClient kubernetes.Interface) error {
					if err := k8sClient.CoreV1().ConfigMaps(namespace).Delete(configMapName, nil); err != nil {
						return microerror.Mask(err)
					}
					return nil
				},
				func(k8sClient kubernetes.Interface) error {
					if _, err := k8sClient.CoreV1().ConfigMaps(namespace).Create(configmap); err != nil {
						return microerror.Mask(err)
					}
					return nil
				},
				func(k8sClient kubernetes.Interface) error {
					if err := k8sClient.CoreV1().ConfigMaps(namespace).Delete(configMapName, nil); err != nil {
						return microerror.Mask(err)
					}
					return nil
				},
			},

			// Check that only one metric is returned, as we only had one resource,
			// despite there being both a create and a delete event.
			expectedMetrics: []*prometheus.Desc{
				informer.CreationTimestampDescription,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			k8sClient, err := newK8sClient()
			if err != nil {
				t.Fatal("expected", nil, "got", err)
			}

			if err := mustSetup(k8sClient); err != nil {
				t.Fatal("expected", nil, "got", err)
			}
			defer mustTeardown(k8sClient)

			operatorkitInformer, err := newOperatorkitInformer(k8sClient)
			if err != nil {
				t.Fatal("expected", nil, "got", err)
			}

			for _, setupEvent := range tc.setupEvents {
				if err := setupEvent(k8sClient); err != nil {
					t.Fatal("expected", nil, "got", err)
				}
			}

			metrics := []prometheus.Metric{}
			metricsChan := make(chan prometheus.Metric)

			go func() {
				operatorkitInformer.Collect(metricsChan)
				close(metricsChan)
			}()

			for _, event := range tc.events {
				if err := event(k8sClient); err != nil {
					t.Fatal("expected", nil, "got", err)
				}
			}

			for metric := range metricsChan {
				metrics = append(metrics, metric)
			}

			// Check that the correct number of metrics were returned.
			if len(metrics) != len(tc.expectedMetrics) {
				descs := []*prometheus.Desc{}
				for _, m := range metrics {
					descs = append(descs, m.Desc())
				}
				t.Fatal("expected", tc.expectedMetrics, "got", descs)
			}

			for i, metric := range metrics {
				t.Logf("saw metric: %v", metric.Desc())

				// Check the correct metric was returned.
				if metric.Desc() != tc.expectedMetrics[i] {
					t.Fatal("expected", tc.expectedMetrics[i], "got", metric.Desc())
				}

				// Check that the timestamp is recent enough.
				modelMetric := prommodel.Metric{}
				if err := metric.Write(&modelMetric); err != nil {
					t.Fatal("expected", nil, "got", err)
				}

				timestamp := int64(*modelMetric.GetGauge().Value)
				now := time.Now().Unix()
				t.Logf("metric value: %v", now)

				if timestamp < now-grace || timestamp > now+grace {
					t.Fatal("expected", fmt.Sprintf("value between %v and %v", now-grace, now+grace), "got", timestamp)
				}
			}
		})
	}
}
