package p2p

import "net"

// Peer is an interface the reprensents the remote node
type Peer interface{
	net.Conn
	Send([]byte) error
	// Close() error
	// RemoteAddress() net.Addr
}


// Transport is anything that handles the communication
// between the nodes in the network. This can be of form(TCP, UDP, Websocket)
type Transport interface{
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
	ListenAddr() string
}