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

const defaultRootFolderName = "BKNet"

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

func (p PathKey) HeadPath() string {
	return strings.Split(p.Pathname, "/")[0]
}

func (p PathKey) Fullpath() string {
	return fmt.Sprintf("%s/%s", p.Pathname, p.Filename)
}

type StoreOptions struct {
	Root          string
	TransformFunc TransformFunc
}

var DefaultPathTransform = func(key string) PathKey {
	return PathKey{
		Pathname: key,
		Filename: key,
	}
}

type Store struct {
	StoreOptions
}

func NewStore(opts StoreOptions) *Store {

	if opts.TransformFunc == nil {
		opts.TransformFunc = DefaultPathTransform
	}

	if opts.Root == "" {
		opts.Root = defaultRootFolderName
	}

	return &Store{
		StoreOptions: opts,
	}
}

func (s *Store) Has(key string) bool {
	PathKey := s.TransformFunc(key)

	fullpathWRoot := fmt.Sprintf("%s/%s", s.Root, PathKey.Fullpath())
	_, err := os.Stat(fullpathWRoot)
	if err != nil {
		return false
	}
	return true
}

func (s *Store) Delete(key string) error {
	pathKey := s.TransformFunc(key)

	defer func() {
		log.Printf("deleted [%s] from disk", pathKey.Filename)
	}()

	firstPathWRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.HeadPath())

	return os.RemoveAll(firstPathWRoot)
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathKey := s.TransformFunc(key)
	pathWRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.Pathname)
	if err := os.MkdirAll(pathWRoot, os.ModePerm); err != nil && !os.IsExist(err) {
		return err
	}
	fullpwroot := fmt.Sprintf("%s/%s", s.Root, pathKey.Fullpath())

	file, err := os.Create(fullpwroot)
	if err != nil {
		return err
	}

	n, err := io.Copy(file, r)
	if err != nil {
		return err
	}

	log.Printf("written (%d) bytes to disk: %s", n, fullpwroot)

	return nil
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.TransformFunc(key)
	fullpwroot := fmt.Sprintf("%s/%s", s.Root, pathKey.Fullpath())
	return os.Open(fullpwroot)
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
