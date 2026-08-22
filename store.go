package main

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strings"
)

func CASTransformFunc(key string) string {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blocksize := 10
	sliceLen := len(hashStr) / blocksize

	paths := make([]string, sliceLen)
	for i := 0; i < sliceLen; i++ {
		from, to := i*blocksize, (i*blocksize)+blocksize
		paths[i] = hashStr[from:to]
	}

	return strings.Join(paths, "/")

}

type TransformFunc func(string) string

type StoreOptions struct {
	TransformFunc TransformFunc
}

var DefaultPathTransform = func(s string) string { return s }

type Store struct {
	StoreOptions
}

func NewStore(opts StoreOptions) *Store {
	return &Store{
		StoreOptions: opts,
	}
}

func (s *Store) write(key string, r io.Reader) error {
	pathname := s.TransformFunc(key)
	if err := os.MkdirAll(pathname, os.ModePerm); err != nil && !os.IsExist(err) {
		return err
	}

	filename := "myfilename"

	pathAndFilename := pathname + "/" + filename

	file, err := os.Create(pathAndFilename)
	if err != nil {
		return err
	}

	n, err := io.Copy(file, r)
	if err != nil {
		return err
	}

	log.Printf("written (%d) bytes to disk: %s", n, pathAndFilename)

	return nil
}
