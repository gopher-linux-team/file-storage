package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func CASTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blocksize := 10
	sliceLen := len(hashStr) / blocksize

	paths := make([]string, sliceLen)
	for i := 0; i < sliceLen; i++ {
		from, to := i*blocksize, (i*blocksize)+blocksize
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		Pathname: strings.Join(paths, "/"),
		Filename: hashStr,
	}

}

type TransformFunc func(string) PathKey

type PathKey struct {
	Pathname string
	Filename string
}

func (p PathKey) Fullpath() string {
	return fmt.Sprintf("%s/%s", p.Pathname, p.Filename)
}

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
	pathKey := s.TransformFunc(key)
	if err := os.MkdirAll(pathKey.Pathname, os.ModePerm); err != nil && !os.IsExist(err) {
		return err
	}

	fullpath := pathKey.Fullpath()

	file, err := os.Create(fullpath)
	if err != nil {
		return err
	}

	n, err := io.Copy(file, r)
	if err != nil {
		return err
	}

	log.Printf("written (%d) bytes to disk: %s", n, fullpath)

	return nil
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.TransformFunc(key)
	return os.Open(pathKey.Fullpath())
}

func (s *Store) Read(key string) (io.Reader, error) {
	file, err := s.readStream(key)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, file)

	return buf, err
}
