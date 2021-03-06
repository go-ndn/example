package main

import (
	"os"
	"time"

	"github.com/go-ndn/mux"
	"github.com/go-ndn/ndn"
	"github.com/go-ndn/packet"
	"github.com/go-ndn/persist"
	"github.com/sirupsen/logrus"
)

func main() {
	// connect to nfd
	conn, err := packet.Dial("tcp", ":6363")
	if err != nil {
		logrus.Error(err)
		return
	}

	// create a new face
	recv := make(chan *ndn.Interest)
	face := ndn.NewFace(conn, recv)
	defer face.Close()

	// read producer key
	pem, err := os.Open("key/default.pri")
	if err != nil {
		logrus.Error(err)
		return
	}
	defer pem.Close()
	key, _ := ndn.DecodePrivateKey(pem)

	// create an interest mux
	m := mux.New()
	// 7. logging
	m.Use(mux.Logger)
	// 6. versioning
	m.Use(mux.Versioner)
	// 5. before encrypting it, zip it
	m.Use(mux.Gzipper)
	// 4. before segmenting it, encrypt it
	m.Use(mux.Encryptor("/producer/encrypt", key.(*ndn.RSAKey)))
	// 3. if the data packet is too large, segment it
	m.Use(mux.Segmentor(10))
	// 2. reply the interest with the on-disk cache
	m.Use(persist.Cacher("test.db"))
	// 1. reply the interest with the in-memory cache
	m.Use(mux.Cacher)
	// 0. an interest packet comes
	m.Use(mux.Queuer)

	// serve encryption key from cache
	m.HandleFunc("/producer/encrypt", func(w ndn.Sender, i *ndn.Interest) error { return nil })

	// serve hello message
	m.HandleFunc("/hello", func(w ndn.Sender, i *ndn.Interest) error {
		return w.SendData(&ndn.Data{
			Name:    ndn.NewName("/hello"),
			Content: []byte(time.Now().UTC().String()),
		})
	})

	// also serve any file under /etc
	m.Handle(mux.FileServer("/file", "/etc"))

	// publish public key file
	m.Handle(mux.StaticFile("key/default.ndncert"))

	// pump the face's incoming interests into the mux
	m.Run(face, recv, key)
}
