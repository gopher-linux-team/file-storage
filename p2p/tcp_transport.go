package p2p

import (
	"fmt"
	"net"
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

// Close implements the Peer interface.
// It closes the underlying TCP connection.
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

type TCPTransportOptions struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TCPTransport struct {
	TCPTransportOptions
	listener net.Listener
	rpcChan  chan RPC
}

func NewTCPTransport(conf TCPTransportOptions) *TCPTransport {
	return &TCPTransport{
		TCPTransportOptions: conf,
		rpcChan:             make(chan RPC),
	}
}

// Consume implements the Transport interface. It returns a read channel.
// This channel is used to receive RPC messages from the transport layer.
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcChan
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
	var err error

	defer func() {
		fmt.Printf("Dropping Peer connection: %s", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, true)

	if err := t.HandshakeFunc(peer); err != nil {
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	// Read loop for incoming RPC messages
	rpc := RPC{}
	for {
		err := t.Decoder.Decode(conn, &rpc)
		if err != nil {

			fmt.Printf("TCP error: %s\n", err)
			return
		}
		rpc.From = conn.RemoteAddr()
		t.rpcChan <- rpc
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
