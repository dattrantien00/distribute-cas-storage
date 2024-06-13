package main

import (
	"cas-storage/p2p"
	"time"
)

func OnPeer(p p2p.Peer) error {
	p.Close()
	return nil
}
func main() {
	tcptransportOpts := p2p.TCPTransportOps{
		ListenAddress: ":3000",
		Decoder:       p2p.DefaultDecoder{},
		HandshakeFunc: p2p.NOPHandshakeFunc,
		// OnPeer:        OnPeer,
	}
	tcpTransport := p2p.NewTCPTransport(tcptransportOpts)
	fileServerOpts := FileServerOpts{
		StorageRoot:       "3000_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
	}

	s := NewFileServer(fileServerOpts)
	go func() {
		time.Sleep(2 * time.Second)
		s.Close()
	}()
	s.Start()

	select{}

	// ts := p2p.NewTCPTransport(opts)
	// if err := ts.ListenAndAccept(); err != nil {
	// 	log.Fatalln(err)
	// }

	// for {
	// 	msg := <-ts.Consume()
	// 	fmt.Println(string(msg.Payload))
	// }
	// select {}
}
