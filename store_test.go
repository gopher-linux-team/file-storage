package main

import (
	"bytes"
	"testing"
)

func TestStore(t *testing.T) {
	opts := StoreOptions{
		TransformFunc: CASTransformFunc,
	}
	store := NewStore(opts)

	data := bytes.NewReader([]byte("somepngbytes"))

	if err := store.write("somepicture", data); err != nil {
		t.Error(err)
	}
}

func TestTransformFunc(t *testing.T) {
	key := "beevswasp"
	pathname := CASTransformFunc(key)
	expected := "5d54b55ab6/78de0138f3/367bff47ae/3c5453de96"
	if pathname != expected {
		t.Errorf("have %s wants %s", pathname, expected)
	}
}
