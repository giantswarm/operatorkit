package informer

import (
	"encoding/json"
	"io"

	"k8s.io/client-go/pkg/runtime"
	"k8s.io/client-go/pkg/watch"
)

type decoder struct {
	stream io.ReadCloser
	obj    ZeroObjectFactory
}

func (d *decoder) Decode() (action watch.EventType, object runtime.Object, err error) {
	var e struct {
		Type   watch.EventType
		Object runtime.Object
	}

	dec := json.NewDecoder(d.stream)
	e.Object = d.obj.NewObject()
	if err := dec.Decode(&e); err != nil {
		return watch.Error, nil, err
	}

	return e.Type, e.Object, nil
}

func (d *decoder) Close() {
	d.stream.Close()
}
