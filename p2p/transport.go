package p2p

import "net"

// Represents a remote peer in the network. This interface is intentionally left empty, as it is meant to be implemented by any type that represents a peer in the network.
type Peer interface {
	Send([]byte) error
	RemoteAddr() net.Addr
	Close() error
}

// Anything that handles communication between peers should implement this interface.
// Any form (UDP, TCP, QUIC, etc.) of transport should implement this interface.
type Transport interface {
	Dial(addr string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
