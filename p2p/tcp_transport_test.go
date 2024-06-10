package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	listenAddr := ":4000"
	ts := NewTCPTransport(listenAddr)
	assert.Equal(t,ts.listenAddress,listenAddr)

	assert.Nil(t,ts.ListenAndAccept())
}
