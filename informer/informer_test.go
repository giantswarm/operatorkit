package informer

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	clientsetfake "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"

	informerfake "github.com/giantswarm/operatorkit/informer/fake"
)

func Test_Informer_Watch(t *testing.T) {
	closer := make(chan struct{}, 1)
	crdClient := clientsetfake.NewSimpleClientset()
	CRD := informerfake.MustNewCRD()

	//err := crd.Ensure(context.TODO(), CRD, crdClient, &backoff.StopBackOff{})
	//if err != nil {
	//	if err != nil {
	//		t.Fatalf("expected %#v got %#v", nil, err)
	//	}
	//}

	go func() {
		e := filepath.Join("apis", "example.com", "v1", "watch", "tests")

		factory := &ZeroObjectFactoryFuncs{
			NewObjectFunc:     func() runtime.Object { return &informerfake.CustomObject{} },
			NewObjectListFunc: func() runtime.Object { return &informerfake.List{} },
		}

		fmt.Printf("%#v\n", crdClient)
		fmt.Printf("%#v\n", crdClient.Discovery())
		fmt.Printf("%#v\n", crdClient.Discovery().RESTClient())

		stream, err := crdClient.Discovery().RESTClient().Get().AbsPath(e).Stream()
		if err != nil {
			t.Fatalf("expected %#v got %#v", nil, err)
		}
		watcher := watch.NewStreamWatcher(newDecoder(stream, factory))

		defer watcher.Stop()

		for {
			select {
			case event, ok := <-watcher.ResultChan():
				fmt.Printf("received event from watcher\n")
				fmt.Printf("%#v\n", ok)
				fmt.Printf("%#v\n", event)
			case <-closer:
				return
			}
		}
	}()

	time.Sleep(1 * time.Second)

	{
		b := informerfake.NewCRO("al7qy")
		e := CRD.ListEndpoint()

		err := crdClient.Discovery().RESTClient().Post().AbsPath(e).Body(b).Do().Error()
		if err != nil {
			t.Fatalf("expected %#v got %#v", nil, err)
		}
	}

	time.Sleep(1 * time.Second)

	{
		b := informerfake.NewCRO("al8qy")
		e := CRD.ListEndpoint()

		err := crdClient.Discovery().RESTClient().Post().AbsPath(e).Body(b).Do().Error()
		if err != nil {
			t.Fatalf("expected %#v got %#v", nil, err)
		}
	}

	time.Sleep(1 * time.Second)

	{
		e := CRD.ResourceEndpoint("al7qy")

		err := crdClient.Discovery().RESTClient().Delete().AbsPath(e).Do().Error()
		if err != nil {
			t.Fatalf("expected %#v got %#v", nil, err)
		}
	}

	time.Sleep(1 * time.Second)

	close(closer)
}
