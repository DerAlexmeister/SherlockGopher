package sherlockanalyser

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	proto "github.com/ob-algdatii-20ss/SherlockGopher/analyser/proto"
	log "github.com/sirupsen/logrus"
)

/*
WebsiteData will receive the website data via grpc.
*/
//nolint: gomnd ineffassign
func (analyser *AnalyserServiceHandler) WebsiteData(ctx context.Context, stream proto.Analyser_WebsiteDataStream) error {
	log.Debug("Receiver->WebsiteData->Receiving website data from crawler")
	buffer := make([][]byte, 0)

	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Error("Receiver->WebsiteData->Stream of request failed")
			return err
		}

		buffer = append(buffer, data.Content)

		if err = stream.Send(&proto.CrawlerAck{Id: data.Id}); err != nil {
			log.Error("Streamer->SendData->Send failed!")
			return err
		}
	}

	// INFO: Drin lassen, wenn wieder error kommt
	if err := stream.Close(); err != nil {
		log.Error("Receiver->WebSiteData->Stream close failed!")
		return err
	}

	data := analyser.InterpretHeaderPackage(buffer[0])
	for i := 1; i < len(buffer); i++ {
		data.cd.responseBodyBytes = append(data.cd.responseBodyBytes, buffer[i]...)
	}

	analyserTaskRequest := NewTask(data.cd)

	analyser.getQueue().AppendQueue(analyserTaskRequest)
	return nil
}

type InterpretedData struct {
	cd              *CrawlerData
	bodyPacketCount int
}

func (analyser *AnalyserServiceHandler) InterpretHeaderPackage(data []byte) *InterpretedData {
	interpretedData := InterpretedData{}
	crawlerData := CrawlerData{}
	interpretedData.cd = &crawlerData

	// Task ID
	start := 0
	end := 8
	addEight := 8
	debug := data[start:end]
	id := binary.LittleEndian.Uint64(debug)
	crawlerData.setTaskID(id)

	// Address Size
	start = end
	end = start + addEight
	debug = data[start:end]
	addressSize := binary.LittleEndian.Uint64(debug)

	// Address
	start = end
	end = start + int(addressSize)
	crawlerData.setAddr(string(data[start:end]))

	// Status Code
	start = end
	end = start + addEight
	debug = data[start:end]
	crawlerData.setStatusCode(int(binary.LittleEndian.Uint64(debug)))

	// Response Time
	start = end
	end = start + addEight

	responseTime := time.Duration(binary.LittleEndian.Uint64(data[start:end]))
	crawlerData.setResponseTime(responseTime)

	// Error String Size
	start = end
	end = start + addEight
	debug = data[start:end]
	errorSize := binary.LittleEndian.Uint64(debug)

	// Error String
	if errorSize != 0 {
		start = end
		end = start + int(errorSize)
		crawlerData.setTaskError(fmt.Errorf(string(data[start:end])))
	} else {
		crawlerData.setTaskError(fmt.Errorf(""))
	}

	// Header Fields Count
	start = end
	end = start + addEight
	headerFieldCount := int(binary.LittleEndian.Uint64(data[start:end]))

	// Header Fields
	headerMap := http.Header{}
	for i := 0; i < headerFieldCount; i++ {
		start = end
		end = start + addEight
		headerFieldSize := int(binary.LittleEndian.Uint64(data[start:end]))

		start = end
		end = start + headerFieldSize
		headerFieldContent := string(data[start:end])
		keyValuePair := strings.Split(headerFieldContent, "W_J")
		headerValues := strings.Split(keyValuePair[1], "J_W")

		for _, headerValue := range headerValues {
			headerMap.Add(keyValuePair[0], headerValue)
		}
	}
	crawlerData.setResponseHeader(&headerMap)

	// Header Fields Count
	start = end
	end = start + addEight
	bodyPacketCount := int(binary.LittleEndian.Uint64(data[start:end]))
	interpretedData.bodyPacketCount = bodyPacketCount

	return &interpretedData
}
