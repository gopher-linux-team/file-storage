package main

import (
	"fmt"
	"log"

	"github.com/gopher-linux-team/file-storage/p2p"
)

type FSrvOpts struct {
	ListenAddr    string
	StorageRoot   string
	PathTransform TransformFunc
	Transport     p2p.Transport
}

type FileServer struct {
	FSrvOpts
	store  *Store
	exitch chan struct{}
}

func NewFileServer(opts FSrvOpts) *FileServer {
	sOpts := StoreOptions{
		Root:          opts.StorageRoot,
		TransformFunc: opts.PathTransform,
	}
	return &FileServer{
		FSrvOpts: opts,
		store:    NewStore(sOpts),
		exitch:   make(chan struct{}),
	}
}

func (s *FileServer) Stop() {
	close(s.exitch)
}

func (s *FileServer) Run() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	s.loop()
	return nil
}

func (s *FileServer) loop() {

	defer func() {
		log.Println("file server stopped, user exit action")
		s.Transport.Close()
	}()

	for {
		select {
		case msg := <-s.Transport.Consume():
			fmt.Println("Received message:", msg)
		case <-s.exitch:
			fmt.Println("Exiting loop")
			return
		}
	}
}
