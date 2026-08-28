package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestStore(t *testing.T) {
	opts := StoreOptions{
		TransformFunc: CASTransformFunc,
	}
	store := NewStore(opts)
	key := "testfilename"
	data := []byte("somepngbytes")

	if err := store.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if ok := store.Has(key); !ok {
		t.Errorf("expected to have key: %s", key)
	}

	r, err := store.Read(key)
	if err != nil {
		t.Error(err)
	}

	b, _ := io.ReadAll(r)
	if !bytes.Equal(b, data) {
		t.Errorf("want %s, have %s", data, b)
	}

	fmt.Println(string(b))

	store.Delete(key)

}

func TestTransformFunc(t *testing.T) {
	key := "beevswasp"
	pathKey := CASTransformFunc(key)
	expectedOriginalKey := "5d54b55ab678de0138f3367bff47ae3c5453de96"
	expectedPath := "5d54b55ab6/78de0138f3/367bff47ae/3c5453de96"
	if pathKey.Pathname != expectedPath {
		t.Errorf("have %s wants %s", pathKey.Pathname, expectedPath)
	}
	if pathKey.Filename != expectedOriginalKey {
		t.Errorf("have %s wants %s", pathKey.Filename, expectedOriginalKey)
	}
}

func TestDelete(t *testing.T) {
	opts := StoreOptions{
		TransformFunc: CASTransformFunc,
	}
	store := NewStore(opts)

	key := "testfilename"
	data := []byte("somepngbytes")

	if err := store.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if err := store.Delete(key); err != nil {
		t.Error(err)
	}

}
