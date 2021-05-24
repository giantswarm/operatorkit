package collector

import (
	"strconv"
	"testing"

	"github.com/giantswarm/k8sclient/v5/pkg/k8sclienttest"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger/microloggertest"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	pkgruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake" //nolint:staticcheck
)

func Test_Timestamp(t *testing.T) {
	pods := []pkgruntime.Object{
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pod-1",
				Namespace: "ns-1",
				Labels: map[string]string{
					"a": "b",
				},
			},
		},
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pod-2",
				Namespace: "ns-2",
				Labels: map[string]string{
					"a": "c",
				},
			},
		},
	}

	testCases := []struct {
		name          string
		objects       []pkgruntime.Object
		selector      labels.Selector
		expectedCount int
	}{
		{
			name:          "case 0: select everything",
			objects:       nil,
			expectedCount: 2,
			selector:      labels.Everything(),
		},
		{
			name:          "case 1: select nothing",
			objects:       nil,
			expectedCount: 0,
			selector:      labels.Nothing(),
		},
		{
			name:          "case 2: select by label, some matches",
			objects:       nil,
			expectedCount: 1,
			selector: labels.SelectorFromSet(labels.Set{
				"a": "b",
			}),
		},
		{
			name:          "case 3: select by label, no matches",
			objects:       nil,
			expectedCount: 0,
			selector: labels.SelectorFromSet(labels.Set{
				"a": "d",
			}),
		},
	}

	for i, tc := range testCases {
		tc := tc
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Log(tc.name)

			clients := k8sclienttest.NewClients(k8sclienttest.ClientsConfig{
				CtrlClient: fake.NewFakeClientWithScheme(scheme.Scheme, pods...),
			})

			config := TimestampConfig{
				Logger:    microloggertest.New(),
				K8sClient: clients,
				NewRuntimeObjectFunc: func() pkgruntime.Object {
					return new(corev1.Pod)
				},
				Selector:   tc.selector,
				Controller: "test",
			}
			collector, err := NewTimestamp(config)
			if err != nil {
				t.Fatal(err)
			}
			metrics := collect(t, collector)

			require.Equal(t, tc.expectedCount, len(metrics))
		})
	}
}

func collect(t *testing.T, collector *Timestamp) []string {
	metrics := make(chan prometheus.Metric)
	done := make(chan bool)
	var output []string

	go func() {
		for {
			m, more := <-metrics
			if more {
				metric := dto.Metric{}
				err := m.Write(&metric)
				if err != nil {
					panic(microerror.JSON(err))
				}
				output = append(output, metric.String())
			} else {
				done <- true
				return
			}
		}
	}()

	go func() {
		err := collector.Collect(metrics)
		if err != nil {
			panic(microerror.JSON(err))
		}
		close(metrics)
	}()

	<-done

	return output
}
