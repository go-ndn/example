package main

import (
	"fmt"
	"net"

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
	face := ndn.NewFace(conn, nil)
	defer face.Close()

	spew.Dump(mux.Assemble(face, ndn.NewName("/hello")))
}
