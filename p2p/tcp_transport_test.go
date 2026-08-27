package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	listenAddr := ":3000"
	tr := NewTCPTransport(TCPTransportOptions{
		ListenAddr:    listenAddr,
		HandshakeFunc: NOPHandshakeFunc,
		Decoder:       DefaultDecoder{},
	})

	assert.Equal(t, tr.ListenAddr, ":3000")

	// Server
	assert.Nil(t, tr.ListenAndAccept())
}
