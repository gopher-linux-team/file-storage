package p2p

// Represents a remote peer in the network. This interface is intentionally left empty, as it is meant to be implemented by any type that represents a peer in the network.
type Peer interface {
}

// Anything that handles communication between peers should implement this interface.
// Any form (UDP, TCP, QUIC, etc.) of transport should implement this interface.
type Transport interface {
	ListenAndAccept() error
}
