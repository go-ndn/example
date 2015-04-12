package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-ndn/mux"
	"github.com/go-ndn/ndn"
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

	var key ndn.Key
	pem, _ := ioutil.ReadFile("key/default.pri")
	key.DecodePrivateKey(pem)

	ndn.Register(face, "/hello", &key)

	m := mux.New()
	m.Handle("/hello", func(face *ndn.Face, i *ndn.Interest) {
		spew.Dump(i)
		face.SendData(&ndn.Data{
			Name:    i.Name,
			Content: []byte(time.Now().UTC().String()),
		})
	})
	m.Run(face, recv)
}
