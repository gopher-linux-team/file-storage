package main

import (
	"github.com/gopher-linux-team/file-storage/p2p"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcpTransportOpts := p2p.TCPTransportOptions{
		ListenAddr:    listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}

	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)
	fileSrvopts := FSrvOpts{
		StorageRoot:      listenAddr + "_BKNet",
		PathTransform:    CASTransformFunc,
		Transport:        tcpTransport,
		EstablishedNodes: nodes,
	}

	s := NewFileServer(fileSrvopts)

	tcpTransport.OnPeer = s.OnPeer

	return s
}

func main() {

	srv1 := makeServer(":3000", "")
	srv2 := makeServer(":4000", ":3000")

	go func() {
		srv1.Run()
	}()

	srv2.Run()

}
