package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/go-ndn/mux"
	"github.com/go-ndn/ndn"
	"github.com/go-ndn/persist"
)

func main() {
	// connect to nfd
	conn, err := net.Dial("tcp", ":6363")
	if err != nil {
		fmt.Println(err)
		return
	}

	// create a new face
	recv := make(chan *ndn.Interest)
	face := ndn.NewFace(conn, recv)
	defer face.Close()

	// read producer key
	pem, err := os.Open("key/default.pri")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pem.Close()
	key, _ := ndn.DecodePrivateKey(pem)

	// a helper to register prefix on nfd
	register := func(name string) {
		ndn.SendControl(face, "rib", "register", &ndn.Parameters{
			Name: ndn.NewName(name),
		}, key)
	}

	register("/hello")
	register("/file")

	// create an interest mux
	m := mux.New()
	// 7. publish public key file
	m.Use(mux.StaticFile("key/default.ndncert"))
	// 6. logging before the interest reaches a handler
	m.Use(mux.Logger)
	// 5. before encrypting it, zip it
	m.Use(mux.Gzipper)
	// 4. before segmenting it, encrypt it
	m.Use(mux.AESEncryptor([]byte("example key 1234")))
	// 3. if the data packet is too large, segment it
	m.Use(mux.Segmentor(10))
	// 2. reply the interest with the on-disk cache
	m.Use(persist.Cacher("test.db"))
	// 1. reply the interest with the in-memory cache
	m.Use(mux.Cacher)
	// 0. an interest packet comes

	// serve hello message
	m.HandleFunc("/hello", func(w ndn.Sender, i *ndn.Interest) {
		w.SendData(&ndn.Data{
			Name:    ndn.NewName("/hello"),
			Content: []byte(time.Now().UTC().String()),
		})
	})

	// also serve any file under /etc
	m.Handle("/file", mux.FileServer("/file", "/etc"))

	// pump the face's incoming interests into the mux
	m.Run(face, recv)
}
