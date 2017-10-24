package main

import (
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-ndn/mux"
	"github.com/go-ndn/ndn"
	"github.com/go-ndn/packet"
	"github.com/go-ndn/tlv"
	"github.com/sirupsen/logrus"
)

func main() {
	// connect to nfd
	conn, err := packet.Dial("tcp", ":6363")
	if err != nil {
		logrus.Error(err)
		return
	}
	// start a new face but do not receive new interests
	face := ndn.NewFace(conn, nil)
	defer face.Close()

	// read producer key
	pem, err := os.Open("../producer/key/default.pri")
	if err != nil {
		logrus.Error(err)
		return
	}
	defer pem.Close()
	key, _ := ndn.DecodePrivateKey(pem)

	// create a data fetcher
	f := mux.NewFetcher()
	// 0. a data packet comes
	// 1. verifiy checksum
	f.Use(mux.ChecksumVerifier)
	// 2. add the data to the in-memory cache
	f.Use(mux.Cacher)
	// 3. logging
	f.Use(mux.Logger)
	// see producer
	dec := mux.Decryptor(key.(*ndn.RSAKey))
	// 4. assemble segments if the content has multiple segments
	// 5. decrypt
	// 6. unzip
	spew.Dump(f.Fetch(face, &ndn.Interest{Name: ndn.NewName("/ndn/guest/alice/1434508942077/KEY/%00%00")}))
	spew.Dump(f.Fetch(face, &ndn.Interest{Name: ndn.NewName("/hello")}, mux.Assembler, dec, mux.Gunzipper))
	spew.Dump(f.Fetch(face, &ndn.Interest{Name: ndn.NewName("/file/hosts")}, mux.Assembler, dec, mux.Gunzipper))

	content, err := f.Fetch(face,
		&ndn.Interest{
			Name: ndn.NewName("/localhop/nfd/rib/list"),
			Selectors: ndn.Selectors{
				MustBeFresh: true,
			},
		}, mux.Assembler)
	if err != nil {
		logrus.Error(err)
		return
	}

	// see nfd
	var rib []ndn.RIBEntry
	tlv.Unmarshal(content, &rib, 128)
	spew.Dump(rib)
}
