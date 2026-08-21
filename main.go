package main

import (
	"fmt"
	"log"

	"github.com/gopher-linux-team/file-storage/p2p"
)

func main() {
	fmt.Println("Running")
	tcpOpts := p2p.TCPTransportOptions{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer: func(peer p2p.Peer) error {
			fmt.Printf("New peer connected: %+v\n", peer)
			return nil
		},
	}
	tr := p2p.NewTCPTransport(tcpOpts)

	go func() {
		for {
			msg := tr.Consume()
			fmt.Printf("%+v\n", msg)
		}
	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatalf("Error starting TCP transport: %v", err)
	}

	fmt.Println("Listening")
	select {}
}
