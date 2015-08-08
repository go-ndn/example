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
	// connect to nfd
	conn, err := net.Dial("tcp", ":6363")
	if err != nil {
		fmt.Println(err)
		return
	}
	// start a new face but do not receive new interests
	face := ndn.NewFace(conn, nil)
	defer face.Close()

	// create a data fetcher
	f := mux.NewFetcher()
	// 0. a data packet comes
	// 1. verifiy checksum
	f.Use(mux.ChecksumVerifier)
	// 2. add the data to the in-memory cache
	f.Use(mux.Cacher)
	// 3. logging
	f.Use(mux.Logger)
	// 4. assemble segments if the content has multiple segments
	f.Use(mux.Assembler)
	// 5. decrypt
	dec := mux.AESDecryptor([]byte("example key 1234"))
	// 6. unzip
	// note: middleware can be both global and local to one handler
	// see producer
	spew.Dump(f.Fetch(face, &ndn.Interest{Name: ndn.NewName("/hello")}, dec, mux.Gunzipper))
	spew.Dump(f.Fetch(face, &ndn.Interest{Name: ndn.NewName("/file/hosts")}, dec, mux.Gunzipper))

	// see nfd
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
