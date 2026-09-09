package p2p

import (
	"errors"
	"fmt"
	"log"
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

// Close implements the Transport interface. It closes the underlying TCP listener.
func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

// ListenAndAccept implements the Transport interface. It starts listening for incoming TCP connections and accepts them.
func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}
	go t.acceptLoop()

	log.Printf("TCP transport listening on %s\n", t.ListenAddr)
	return nil
}

type Temp struct{}

// handleConnection handles an incoming TCP connection. It performs the handshake,
// calls the OnPeer callback, and starts reading RPC messages from the connection.
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

// acceptLoop continuously accepts incoming TCP connections and handles them in separate goroutines.
func (t *TCPTransport) acceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			log.Println("TCP listener closed, stopping accept loop")
			return
		}
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
			continue
		}

		go t.handleConnection(conn)
	}
}
