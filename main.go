package main

import (
	"bytes"
	"cas-storage/p2p"
	"log"
	"time"
)

func makeServer(listenAddr, root string, nodes ...string) *FileServer {
	tcptransportOpts := p2p.TCPTransportOps{
		ListenAddress: listenAddr,
		Decoder:       p2p.DefaultDecoder{},
		HandshakeFunc: p2p.NOPHandshakeFunc,
		// OnPeer:        OnPeer,
	}
	tcpTransport := p2p.NewTCPTransport(tcptransportOpts)
	s := NewFileServer(FileServerOpts{
		StorageRoot:       root,
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	})
	tcpTransport.OnPeer = s.OnPeer

	return s
}
func main() {

	s1 := makeServer(":3001", "3000_network")

	go func() {
		log.Fatal(s1.Start())
	}()
	s2 := makeServer(":4001", "4000_network", ":3001")

	go s2.Start()

	time.Sleep(1 * time.Second)

	// s2.peers["127.0.0.1:3000"].Send([]byte("abc"))
	data := bytes.NewBuffer([]byte("hi"))
	s2.StoreData("key",data)

	select {}
}
