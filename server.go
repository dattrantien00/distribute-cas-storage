package main

import (
	"bytes"
	"cas-storage/p2p"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"
	"time"
)

type FileServerOpts struct {
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
	BootstrapNodes    []string
}

type FileServer struct {
	FileServerOpts

	peerLock sync.Mutex
	peers    map[string]p2p.Peer
	store    *Store
	quitch   chan struct{}
}

type Message struct {
	From    string
	Payload any
}

func NewFileServer(opts FileServerOpts) *FileServer {
	return &FileServer{
		FileServerOpts: opts,
		store: NewStore(StoreOpts{
			Root:              opts.StorageRoot,
			PathTransformFunc: opts.PathTransformFunc,
		}),
		quitch: make(chan struct{}),
		// peerLock: sync.Mutex{},
		peers: map[string]p2p.Peer{},
	}

}

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	s.bootstrapNetwork()

	s.loop()
	return nil
}

// store file and broadcast to all network
func (s *FileServer) StoreData(key string, r io.Reader) error {

	// buf := new(bytes.Buffer)
	// tee := io.TeeReader(r, buf)

	// err := s.Store(key, tee)
	// if err != nil {
	// 	return err
	// }

	// payload := &DataMessage{
	// 	Key:  key,
	// 	Data: buf.Bytes(),
	// }
	// fmt.Println(string(buf.Bytes()))

	// return s.broadcast(&Message{
	// 	From:    s.Transport.ListenAddr(),
	// 	Payload: payload,
	// })
	buf := new(bytes.Buffer)
	msg := Message{
		Payload: []byte("hi"),
	}
	if err := gob.NewEncoder(buf).Encode(msg); err != nil {
		return err
	}
	for _, peer := range s.peers {
		peer.Send(buf.Bytes())
	}

	time.Sleep(2*time.Second)
	payload := []byte("This file")
	for _, peer := range s.peers {
		peer.Send(payload)
	}
	return nil
}

func (s *FileServer) broadcast(msg *Message) error {
	// return s.store.Write(key, r)

	peers := []io.Writer{}
	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...)
	return gob.NewEncoder(mw).Encode(msg)
}
func (s *FileServer) Store(key string, r io.Reader) error {
	return s.store.Write(key, r)
}

// func (s *FileServer) handleMessage(message *Message) error {
// 	// switch v:= message.Payload.(type){
// 	// case *DataMessage:
// 	// 	fmt.Println("received message %+v\n",v)

//		// }
//		return nil
//	}
func (s *FileServer) loop() {
	defer func() {
		log.Println("file server stopped due to user quit action")
		s.Transport.Close()
	}()
	for {
		select {
		case rpc := <-s.Transport.Consume():
			
			var msg Message
			if err := gob.NewDecoder(bytes.NewBuffer(rpc.Payload)).Decode(&msg); err != nil {
				log.Println(err)
			}
			fmt.Println("%+v\n", string(msg.Payload.([]byte)))

			peer, ok := s.peers[rpc.From]
			if !ok {
				panic("peer not found in peer map")
			}
			b := make([]byte, 1000)
			if _,err := peer.Read(b);err != nil{
				panic(err)
			}
			fmt.Println(string(b))
			panic("123123123")

			// err := s.handleMessage(&data)
			// if err != nil {
			// 	log.Println(err)
			// }
			// fmt.Printf("%+v\n",(data))
			// fmt.Println("msg: ", string(msg.Payload))
		case <-s.quitch:
			return
		}
	}
}

func (s *FileServer) bootstrapNetwork() error {
	for _, addr := range s.BootstrapNodes {

		if len(addr) == 0 {
			continue
		}
		go func(addr string) {
			fmt.Println("attemping to connect with remote:", addr)
			if err := s.Transport.Dial(addr); err != nil {
				log.Println("dial err:", err)
			}

		}(addr)
	}
	return nil
}

func (s *FileServer) Close() {
	close(s.quitch)
}

func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	s.peers[p.RemoteAddr().String()] = p

	log.Printf("connectted with remote %s", p.RemoteAddr().String())
	return nil
}
