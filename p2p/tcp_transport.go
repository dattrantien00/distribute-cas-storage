package p2p

import (
	"fmt"
	"net"
	"sync"
)

type TCPPeer struct {
	conn net.Conn

	// if we dial and retrieve a conn => outbound = true
	// if we accept a conn => outbound = false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransport struct {
	listenAddress string
	listener      net.Listener
	handshakeFunc HandshakeFunc
	// decoder Decoder

	mu            sync.RWMutex
	peer          map[net.Addr]Peer
}

func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		handshakeFunc:     NOPHandshakeFunc,
		listenAddress: listenAddr,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	ln, err := net.Listen("tcp", t.listenAddress)
	if err != nil {
		return err
	}

	t.listener = ln

	go t.acceptLoop()
	return err
}

func (t *TCPTransport) acceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %s", err)
		}
		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)
	if err := t.handshakeFunc(conn); err != nil{

	}
	fmt.Printf("new incoming connection %+v\n", peer)
}
