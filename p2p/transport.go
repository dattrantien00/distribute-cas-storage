package p2p

// Peer is an interface the reprensents the remote node
type Peer interface{
	Close() error
}


// Transport is anything that handles the communication
// between the nodes in the network. This can be of form(TCP, UDP, Websocket)
type Transport interface{
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}