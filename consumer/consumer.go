package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-ndn/ndn"
)

func segmentFetcher(face *ndn.Face, s string) (content []byte) {
	var start int
	segNum := make([]byte, 8)
	base := ndn.NewName(s)
	for {
		binary.BigEndian.PutUint64(segNum, uint64(start))
		d, ok := <-face.SendInterest(&ndn.Interest{
			Name: ndn.Name{Components: append(base.Components, segNum)},
		})
		if !ok {
			return
		}
		content = append(content, d.Content...)
		if len(d.Name.Components) > 0 &&
			!bytes.Equal(d.Name.Components[len(d.Name.Components)-1], d.MetaInfo.FinalBlockID.Component) {
			start += len(d.Content)
		} else {
			break
		}
	}
	return
}

func main() {
	conn, err := net.Dial("tcp", ":6363")
	if err != nil {
		fmt.Println(err)
		return
	}
	face := ndn.NewFace(conn, nil)
	defer face.Close()

	spew.Dump(segmentFetcher(face, "/hello"))
}
