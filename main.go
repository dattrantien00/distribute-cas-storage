package main

import (
	"cas-storage/p2p"
	"log"
)

func main() {
	ts := p2p.NewTCPTransport(":3000")
	if err := ts.ListenAndAccept(); err != nil {
		log.Fatalln(err)
	}
	select {}
}
