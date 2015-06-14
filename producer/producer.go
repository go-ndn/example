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
	conn, err := net.Dial("tcp", ":6363")
	if err != nil {
		fmt.Println(err)
		return
	}
	recv := make(chan *ndn.Interest)
	face := ndn.NewFace(conn, recv)
	defer face.Close()

	pem, err := os.Open("key/default.pri")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pem.Close()
	key, _ := ndn.DecodePrivateKey(pem)

	register := func(name string) {
		ndn.SendControl(face, "rib", "register", &ndn.Parameters{
			Name: ndn.NewName(name),
		}, key)
	}

	register("/hello")
	register("/file")

	m := mux.New()
	m.Use(mux.Logger)
	m.Use(mux.Gzipper)
	m.Use(mux.AESEncryptor([]byte("example key 1234")))
	m.Use(mux.Segmentor(10))
	m.Use(persist.Cacher("test.db"))
	m.Use(mux.Cacher)
	m.HandleFunc("/hello", func(w ndn.Sender, i *ndn.Interest) {
		w.SendData(&ndn.Data{
			Name:    ndn.NewName("/hello"),
			Content: []byte(time.Now().UTC().String()),
		})
	})
	m.Handle("/file", mux.FileServer("/file", "/etc"))
	m.Run(face, recv)
}
