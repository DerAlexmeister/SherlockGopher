package crawler

import (
	"github.com/golang/protobuf/ptypes/empty"
	proto "github.com/ob-algdatii-20ss/leistungsnachweis-dievierausrufezeichen/crawler/proto"
)

type CrawlerServiceHandler struct{}

const chunkSize = 64 * 1024 // 64 KiB

type chunkerSrv []byte

func NewCrawlerServiceHandler() *CrawlerServiceHandler {
	return &CrawlerServiceHandler{}
}

func (c chunkerSrv) SendHttpPage(_ *empty.Empty, srv proto.Chunker_ChunkerStream) error {
	chnk := &proto.Chunk{}
	for currentByte := 0; currentByte < len(c); currentByte += chunkSize {
		if currentByte+chunkSize > len(c) {
			chnk.Chunk = c[currentByte:len(c)]
		} else {
			chnk.Chunk = c[currentByte : currentByte+chunkSize]
		}
		if err := srv.Send(chnk); err != nil {
			return err
		}
	}
	return nil
}
