// +build integration

package informer

import (
	"context"
	"testing"
	"time"
)

// Test_Informer_Integration_Basic is a integration test for basic cache
// informer operations. The test verifies the informer is operating as expected
// when processing basic sequences of creating and deleting runtime objects.
func Test_Informer_Integration_Basic(t *testing.T) {
	testSetup(t)
	defer testTeardown(t)

	const timeDelta = time.Millisecond * 100

	objectIDOne := "al7qy"
	objectIDTwo := "al8qy"

	ctx, cancelFunc := context.WithCancel(context.Background())
	newInformer := testNewInformer(t, time.Second*2, time.Second*10)

	// We create a custom object before starting the informer watch. This causes
	// the informer to fill the cache and to initially sent cached events to the
	// delete and update channels provided by the watch.
	testCreateObj(t, objectIDOne)

	// When there is a runtime object in the API we start the watch.
	deleteChan, updateChan, errChan := newInformer.Watch(ctx)

	// We define a general control goroutine to stop test execution on errors or
	// after timeouts. The timeout is 25 seconds because of two resync periods,
	// one rate limit wait period, plus some time buffer for all the test magic
	// happening here: (2 * 10) + 2 + buffer.
	go func() {
		for {
			select {
			case <-time.After(25 * time.Second):
				t.Fatalf("expected proper test execution got timeout")
			case err := <-errChan:
				if err != nil {
					t.Fatalf("expected %#v got %#v", nil, err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// This is the first time we want to catch an event from the informer watch.
	// This should give us the runtime object we created before starting the
	// watch.
	{
		select {
		case <-deleteChan:
			t.Fatalf("expected update event got delete event")
		case e := <-updateChan:
			testAssertCROWithID(t, e, objectIDOne)
		}
	}

	// We create another runtime object. This should be received immediately.
	{
		testCreateObj(t, objectIDTwo)

		start := time.Now()

		select {
		case <-deleteChan:
			t.Fatalf("expected update event got delete event")
		case e := <-updateChan:
			testAssertCROWithID(t, e, objectIDTwo)
		}

		d := time.Since(start)
		if !durationEquals(0, d, timeDelta) {
			t.Fatalf("expected %#v got %#v", "round about 0 seconds", d.Seconds())
		}
	}

	// Now nothing happened to the runtime objects within the Kubernetes API and
	// after a while the cache informer does a resync. That means we receive the
	// very first runtime object we created after round about 10 seconds, because
	// the cache informer is configured to resync after 10 seconds.
	{
		start := time.Now()

		select {
		case <-deleteChan:
			t.Fatalf("expected update event got delete event")
		case e := <-updateChan:
			testAssertCROWithID(t, e, objectIDOne, objectIDTwo)
		}

		d := time.Since(start)
		if !durationEquals(10*time.Second, d, timeDelta) {
			t.Fatalf("expected %#v got %#v", "round about 10 seconds", d.Seconds())
		}
	}

	// There are two runtime objects in the Kubernetes API. When the cache
	// informer is resyncing we receive the second runtime object after another 2
	// seconds because the cache informer is configured to rate limit between the
	// runtime objects for 2 seconds.
	{
		start := time.Now()

		select {
		case <-deleteChan:
			t.Fatalf("expected update event got delete event")
		case e := <-updateChan:
			testAssertCROWithID(t, e, objectIDOne, objectIDTwo)
		}

		d := time.Since(start)
		if !durationEquals(2*time.Second, d, timeDelta) {
			t.Fatalf("expected %#v got %#v", "round about 2 seconds", d.Seconds())
		}
	}

	// Now we delete a runtime object. This event is expected to be received
	// immediately.
	{
		testDeleteObj(t, objectIDOne)

		start := time.Now()

		select {
		case e := <-deleteChan:
			testAssertCROWithID(t, e, objectIDOne)
		case <-updateChan:
			t.Fatalf("expected delete event got update event")
		}

		d := time.Since(start)
		if !durationEquals(0, d, timeDelta) {
			t.Fatalf("expected %#v got %#v", "round about 0 seconds", d.Seconds())
		}
	}

	// After another 10 seconds the cache informer does a resync period again and
	// we receive the runtime object left.
	{
		start := time.Now()

		select {
		case <-deleteChan:
			t.Fatalf("expected update event got delete event")
		case e := <-updateChan:
			testAssertCROWithID(t, e, objectIDTwo)
		}

		d := time.Since(start)
		if !durationEquals(10*time.Second, d, timeDelta) {
			t.Fatalf("expected %#v got %#v", "round about 10 seconds", d.Seconds())
		}
	}

	cancelFunc()
}

func durationEquals(expected, actual, delta time.Duration) bool {
	var diff time.Duration
	if expected > actual {
		diff = expected - actual
	} else {
		diff = actual - expected
	}

	return diff < delta
}
