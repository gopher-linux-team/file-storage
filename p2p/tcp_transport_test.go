package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	listenAddr := ":8080"
	tr := NewTCPTransport(TCPTransportOptions{
		ListenAddr:    listenAddr,
		HandshakeFunc: NOPHandshakeFunc,
		Decoder:       GOBDecoder{},
	})

	// Server
	assert.Nil(t, tr.ListenAndAccept())
}
