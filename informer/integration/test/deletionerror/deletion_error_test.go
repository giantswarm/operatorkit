// +build k8srequired

package deletionerror

import (
	"context"
	"testing"
	"time"
)

func Test_Informer_Integration_DeletionError(t *testing.T) {
	mustSetup()
	defer mustTeardown()

	idOne := "al7qy"
	idTwo := "al8qy"
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	operatorkitInformer, filterWatcher, err := newOperatorkitInformerAndWatcher(time.Second*2, time.Second*10)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	// We create two config maps before starting the informer watch. This causes
	// the informer to fill the cache and to initially sent cached events to the
	// delete and update channels provided by the watch.
	err = createConfigMap(idOne)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
	err = createConfigMap(idTwo)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	// When there are the config maps in the Kubernetes API we start the watch.
	deleteChan, updateChan, errChan := operatorkitInformer.Watch(ctx)

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

	// Due to the two created config maps we get two events in the informer which
	// we simply drain here because they stand in our way of testing the actual
	// matter of delete events. When anything else is weird with the informer and
	// the channels hang on this stage then the test will time out due to the
	// safety goroutine above.
	{
		t.Logf("draining update channel")

		<-updateChan
		<-updateChan

		t.Logf("drained update channel")
	}

	// We disable the event dispatching in the filter watcher immediately before
	// any config maps are created in the Kubernetes API.
	{
		t.Logf("disabling informer event dispatching")

		filterWatcher.SetDispatchEvents(false)

		t.Logf("disabled informer event dispatching")
	}

	// Now we delete the first config map. This event is expected to be received
	// immediately in normal circumstances. For this test though, we disabled
	// event dispatching in the filter watcher. Thus we check if there was any
	// event at all. Note that we sleep a bit after deleting the config map to let
	// the delete event go through the whole stack. This may takes a while so we
	// simply set it to 5 seconds.
	{
		err := deleteConfigMap(idOne)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		time.Sleep(5 * time.Second)

		select {
		case <-deleteChan:
			t.Fatalf("did not expect any event but got delete event")
		case <-updateChan:
			t.Fatalf("did not expect any event but got update event")
		case <-time.After(time.Second):
			// fall through
		}
	}

	// Now the first config map got deleted and we ignored the delete event. The
	// informer did not have any chance to perceive the deletion on its own. We
	// enable the event dispatching again to check what config maps the informer
	// still knows about.
	{
		t.Logf("enabling informer event dispatching")

		filterWatcher.SetDispatchEvents(true)

		t.Logf("enabled informer event dispatching")
	}

	// After enabling the event dispatching in the filter watcher again, we should
	// start to receive events again from the informer. Since the first config map
	// got deleted in the Kubernetes API, the only event we should receive is an
	// update event for the second config map. Note that we loop over the event
	// channels because of a potential bug that could make the test succeed
	// occasionally. In case the filter watcher is buggy and dispatches events
	// even though we do not want or expect that, there would be two config maps
	// dispatched through the informer. The order of these two objects would never
	// be guaranteed and we might see the first config map as the one we actually
	// expect while the other that should already be deleted could be queued at
	// the end of the updater channel. To remedy such a situation we loop over the
	// informer channels and ensure we always ever see the second config map which
	// indicates the test setup above worked properly and the test itself does not
	// generate false positives. The loop contains a timeout channel itself in
	// order to end the test after the updater channel got drained completely
	// regardless of itself having contained one or two config maps.
	{
		for {
			select {
			case <-deleteChan:
				t.Fatalf("expected update event got delete event")
			case e := <-updateChan:
				mustAssertWithIDs(e, idTwo)
			case <-time.After(5 * time.Second):
				return
			}
		}
	}
}
