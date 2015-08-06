package main

import (
	"fmt"
	"net"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-ndn/mux"
	"github.com/go-ndn/ndn"
	"github.com/go-ndn/tlv"
)

func main() {
	conn, err := net.Dial("tcp", ":6363")
	if err != nil {
		fmt.Println(err)
		return
	}
	face := ndn.NewFace(conn, nil)
	defer face.Close()

	f := mux.NewFetcher()
	f.Use(mux.BasicVerifier)
	f.Use(mux.Cacher)
	f.Use(mux.Logger)
	f.Use(mux.Assembler)
	dec := mux.AESDecryptor([]byte("example key 1234"))
	spew.Dump(f.Fetch(face, &ndn.Interest{Name: ndn.NewName("/hello")}, dec, mux.Gunzipper))
	spew.Dump(f.Fetch(face, &ndn.Interest{Name: ndn.NewName("/file/hosts")}, dec, mux.Gunzipper))

	var rib []ndn.RIBEntry
	tlv.UnmarshalByte(f.Fetch(face,
		&ndn.Interest{
			Name: ndn.NewName("/localhop/nfd/rib/list"),
			Selectors: ndn.Selectors{
				MustBeFresh: true,
			},
		}),
		&rib,
		128,
	)
	spew.Dump(rib)
}
