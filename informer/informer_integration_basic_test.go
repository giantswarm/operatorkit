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

	ctx, cancelFunc := context.WithCancel(context.Background())
	newInformer := testNewInformer(t, time.Second*2, time.Second*10)

	// We create a custom object before starting the informer watch. This causes
	// the informer to fill the cache and to initially sent cached events to the
	// delete and update channels provided by the watch.
	testCreateCRO(t, "al7qy")

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
				t.Fatalf("expected %#v got %#v", "proper test execution", "timeout")
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
			t.Fatalf("expected %#v got %#v", "update event", "delete event")
		case e := <-updateChan:
			testAssertCROWithID(t, e, "al7qy")
		}
	}

	// We create another runtime object. This should be received immediately.
	{
		testCreateCRO(t, "al8qy")

		start := time.Now()

		select {
		case <-deleteChan:
			t.Fatalf("expected %#v got %#v", "update event", "delete event")
		case e := <-updateChan:
			testAssertCROWithID(t, e, "al8qy")
		}

		d := time.Now().Sub(start)
		if d.Seconds() > 0.1 {
			t.Fatalf("expected %#v got %#v", "round about 0.1 second", d.Seconds())
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
			t.Fatalf("expected %#v got %#v", "update event", "delete event")
		case e := <-updateChan:
			testAssertCROWithID(t, e, "al7qy")
		}

		d := time.Now().Sub(start)
		if d.Seconds() > 10.1 {
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
			t.Fatalf("expected %#v got %#v", "update event", "delete event")
		case e := <-updateChan:
			testAssertCROWithID(t, e, "al8qy")
		}

		d := time.Now().Sub(start)
		if d.Seconds() > 2.1 {
			t.Fatalf("expected %#v got %#v", "round about 2 seconds", d.Seconds())
		}
	}

	// Now we delete a runtime object. This event is expected to be received
	// immediately.
	{
		testDeleteCRO(t, "al7qy")

		start := time.Now()

		select {
		case e := <-deleteChan:
			testAssertCROWithID(t, e, "al7qy")
		case <-updateChan:
			t.Fatalf("expected %#v got %#v", "delete event", "update event")
		}

		d := time.Now().Sub(start)
		if d.Seconds() > 0.1 {
			t.Fatalf("expected %#v got %#v", "round about 0.1 second", d.Seconds())
		}
	}

	// After another 10 seconds the cache informer does a resync period again and
	// we receive the runtime object left.
	{
		start := time.Now()

		select {
		case <-deleteChan:
			t.Fatalf("expected %#v got %#v", "update event", "delete event")
		case e := <-updateChan:
			testAssertCROWithID(t, e, "al8qy")
		}

		d := time.Now().Sub(start)
		if d.Seconds() > 10.1 {
			t.Fatalf("expected %#v got %#v", "round about 10 seconds", d.Seconds())
		}
	}

	cancelFunc()
}
