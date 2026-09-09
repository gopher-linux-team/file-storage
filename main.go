package main

import (
	"log"
	"time"

	"github.com/gopher-linux-team/file-storage/p2p"
)

func main() {

	tcpTransportOpts := p2p.TCPTransportOptions{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		//TODO: onPeer: func(p p2p.Peer) error {}
	}

	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)
	fileSrvopts := FSrvOpts{
		StorageRoot:   "3000_BKNet",
		PathTransform: CASTransformFunc,
		Transport:     tcpTransport,
	}
	srv := NewFileServer(fileSrvopts)

	go func() {
		time.Sleep(time.Second * 3)
		srv.Stop()
	}()

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}

}
