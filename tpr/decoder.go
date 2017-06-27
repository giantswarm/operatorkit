package tpr

import (
	"encoding/json"
	"io"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
)

type decoder struct {
	stream  io.ReadCloser
	obj     ZeroObjectFactory
	jsonDec *json.Decoder
}

func newDecoder(stream io.ReadCloser, obj ZeroObjectFactory) *decoder {
	return &decoder{
		stream:  stream,
		obj:     obj,
		jsonDec: json.NewDecoder(stream),
	}
}

func (d *decoder) Decode() (action watch.EventType, object runtime.Object, err error) {
	var e struct {
		Type   watch.EventType
		Object runtime.Object
	}

	e.Object = d.obj.NewObject()
	if err := d.jsonDec.Decode(&e); err != nil {
		return watch.Error, nil, err
	}

	return e.Type, e.Object, nil
}

func (d *decoder) Close() {
	d.stream.Close()
}
