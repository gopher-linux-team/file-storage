package p2p

import (
	"fmt"
	"net"
	"sync"
)

// Represents the remote node over TCP stablished connection. It implements the Peer interface.
type TCPPeer struct {
	conn     net.Conn
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransportOptions struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
}

type TCPTransport struct {
	TCPTransportOptions
	listener net.Listener

	mux   sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(conf TCPTransportOptions) *TCPTransport {
	return &TCPTransport{
		TCPTransportOptions: conf,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}
	go t.acceptLoop()
	return nil
}

type Temp struct{}

func (t *TCPTransport) handleConnection(conn net.Conn) {
	peer := NewTCPPeer(conn, true)

	if err := t.HandshakeFunc(peer); err != nil {
		fmt.Printf("TCP handshake error: %s\n", err)
		conn.Close()
		return
	}

	lenDecodeError := 0
	// Read loop
	msg := &Temp{}
	for {
		if err := t.Decoder.Decode(conn, msg); err != nil {
			lenDecodeError++
			if lenDecodeError > 5 {
				fmt.Printf("TCP error: %s\n", err)
				return
			}
		}
	}
}

func (t *TCPTransport) acceptLoop() error {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
			continue
		}

		go t.handleConnection(conn)
	}
}
