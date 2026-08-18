package p2p

// HandshakeFunc defines the function signature for performing a handshake with a peer.
type HandshakeFunc func(Peer) error

func NOPHandshakeFunc(Peer) error { return nil }
