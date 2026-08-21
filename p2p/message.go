package p2p

import "net"

// Represents any data exchanged between peers (over transport) in the P2P network.
type RPC struct {
	From    net.Addr
	Payload []byte
}
