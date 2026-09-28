package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/gopher-linux-team/file-storage/p2p"
)

type FSrvOpts struct {
	ListenAddr       string
	StorageRoot      string
	PathTransform    TransformFunc
	Transport        p2p.Transport
	EstablishedNodes []string
}

type FileServer struct {
	FSrvOpts

	store  *Store
	exitch chan struct{}

	peers    map[string]p2p.Peer
	peerLock sync.Mutex
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
		peers:    make(map[string]p2p.Peer),
	}
}

func (s *FileServer) Stop() {
	close(s.exitch)
}

func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	s.peers[p.RemoteAddr().String()] = p
	log.Printf("new peer connected: %s", p.RemoteAddr().String())
	return nil
}

func (s *FileServer) Run() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}
	if len(s.EstablishedNodes) != 0 {
		s.establishNetwork()
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

func (s *FileServer) establishNetwork() error {

	for _, addr := range s.EstablishedNodes {
		if len(addr) == 0 {
			continue
		}

		go func(addr string) {
			fmt.Println("Attempting to connect to established node:", addr)
			if err := s.Transport.Dial(addr); err != nil {
				log.Printf("Failed to connect to %s: %v", addr, err)
			}
		}(addr)
	}

	return nil
}
