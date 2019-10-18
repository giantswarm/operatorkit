package controller

import (
	"testing"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

func Test_HashEvent(t *testing.T) {
	testCases := []struct {
		Event        watch.Event
		Concurrency  int
		ExpectedHash int
	}{
		// Test 0 ensures the hash for a generic object.
		{
			Event: watch.Event{
				Type: watch.Added,
				Object: &v1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foobar",
						Namespace: "blah",
					},
				},
			},
			Concurrency:  32,
			ExpectedHash: 2,
		},

		// Test 1 ensures the hash for another generic object.
		{
			Event: watch.Event{
				Type: watch.Added,
				Object: &v1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "qweqweqwe",
						Namespace: "bloop",
					},
				},
			},
			Concurrency:  32,
			ExpectedHash: 9,
		},

		// Test 2 ensures the hash for yet another generic object.
		{
			Event: watch.Event{
				Type: watch.Added,
				Object: &v1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "dajboubfob",
						Namespace: "foiniof",
					},
				},
			},
			Concurrency:  32,
			ExpectedHash: 22,
		},
	}

	for i, tc := range testCases {
		hash := hashEvent(tc.Event, tc.Concurrency)
		if tc.ExpectedHash != hash {
			t.Fatal("test", i, "expected", tc.ExpectedHash, "got", hash)
		}
	}
}
