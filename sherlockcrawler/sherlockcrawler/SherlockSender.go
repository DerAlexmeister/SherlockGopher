package sherlockcrawler

import (
	"context"
	"encoding/binary"
	"fmt"
	"strings"
	"sync"

	proto "github.com/ob-algdatii-20ss/SherlockGopher/analyser/proto"
	log "github.com/sirupsen/logrus"
)

const (
	packetSize = 307200000
)

//TODO JW Kommentare und logging.

func (sherlock *SherlockCrawler) SendWebsiteData(sender proto.AnalyserService, task *CrawlerTaskRequest, wg *sync.WaitGroup) error {
	if task.GetTaskError() != nil {
		wg.Done()
	} else {
		(*task).setTaskState(PROCESSING)
		dataToSend := sherlock.PrepareData(task)
		mes := fmt.Sprintf("SherlockSender->SendWebsiteData->Start: %v", task.GetAddr())
		log.Info(mes)
		err := sherlock.SendData(task.GetTaskID(), dataToSend, sender, task)
		if err != nil {
			log.Fatal("SendData failed")
		}
		mes = fmt.Sprintf("SherlockSender->SendWebsiteData->End: %v", task.GetAddr())
		log.Info(mes)
		wg.Done()
	}
	return nil
}

/*
SendData sends the website data to the analyser via grpc streams. Receives an ack if everything went right.
*/
func (sherlock *SherlockCrawler) SendData(taskID uint64, data [][]byte, sender proto.AnalyserService, task *CrawlerTaskRequest) error {
	sherlock.Queue.mutex.Lock()

	stream, err := (*sherlock.analyserService).WebsiteData(context.TODO())
	if err != nil {
		log.Error("Streamer->SendData->Stream returned by WebsiteData failed")
		return err
	}

	for i := 0; i < len(data); i++ {
		if err := stream.Send(&proto.CrawlerPackage{
			Id:      taskID,
			Content: data[i],
		}); err != nil {
			log.WithFields(log.Fields{
				"err":  err,
				"task": task.GetAddr(),
			}).Error("Streamer->SendData->Send failed!")
			return err
		}

		rsp, err := stream.Recv()
		if err != nil {
			log.WithFields(log.Fields{
				"err":  err,
				"task": task.GetAddr(),
			}).Error("Streamer->SendData->Receive failed!")
			return err
		}

		if rsp.Id != taskID {
			log.Error("Streamer->SendData->Unknown ID was returned")
		}
	}

	stream.Close()

	sherlock.Queue.mutex.Unlock()
	return nil
}

/*
PrepareData prepares the the website data for the stream. CrawlerTaskRequest is converted to [][]byte.
*/
func (sherlock *SherlockCrawler) PrepareData(task *CrawlerTaskRequest) [][]byte {
	buffer := make([][]byte, 0)
	headerPackage := make([]byte, 0)

	// TaskID
	taskIDBuffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(taskIDBuffer, task.GetTaskID())
	headerPackage = append(headerPackage, taskIDBuffer...)

	// Address
	byteAddress := []byte(task.GetAddr())
	addressSize := uint64(len(byteAddress))
	addressSizeBuffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(addressSizeBuffer, addressSize)
	headerPackage = append(headerPackage, addressSizeBuffer...)
	headerPackage = append(headerPackage, byteAddress...)

	// Status Code
	statusCodeBuffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(statusCodeBuffer, uint64(task.GetStatusCode()))
	headerPackage = append(headerPackage, statusCodeBuffer...)

	// Response Time
	responseTimeBuffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(responseTimeBuffer, uint64(task.GetResponseTime().Nanoseconds()))
	headerPackage = append(headerPackage, responseTimeBuffer...)

	// Error
	var errorSize uint64
	errorSizeBuffer := make([]byte, 8)
	if task.GetTaskError() != nil {
		byteError := []byte(task.GetTaskError().Error())
		errorSize = uint64(len(byteAddress))
		binary.LittleEndian.PutUint64(addressSizeBuffer, errorSize)
		headerPackage = append(headerPackage, errorSizeBuffer...)
		headerPackage = append(headerPackage, byteError...)
	} else {
		binary.LittleEndian.PutUint64(errorSizeBuffer, 0)
		headerPackage = append(headerPackage, errorSizeBuffer...)
	}

	// Header Fields Count
	taskHeaderFieldsCount := make([]byte, 8)
	binary.LittleEndian.PutUint64(taskHeaderFieldsCount, uint64(len(task.GetResponseHeader())))
	headerPackage = append(headerPackage, taskHeaderFieldsCount...)

	// Header Fields
	for k, v := range task.GetResponseHeader() {
		// Header Field Size
		headerValues := strings.Join(v, "J_W")
		headerData := k + "W_J" + headerValues
		headerFieldSize := make([]byte, 8)

		binary.LittleEndian.PutUint64(headerFieldSize, uint64(len(headerData)))
		headerPackage = append(headerPackage, headerFieldSize...)
		headerPackage = append(headerPackage, []byte(headerData)...)
	}

	// Body package count
	taskRespondBody := task.GetResponseBody()
	packetCount := len(taskRespondBody) / packetSize
	lastPacketSize := len(taskRespondBody) % packetSize
	lastPacket := lastPacketSize != 0
	finalPacketCount := packetCount
	if lastPacket {
		finalPacketCount++
	}

	bodyPacketCount := make([]byte, 8)
	binary.LittleEndian.PutUint64(bodyPacketCount, uint64(finalPacketCount))
	headerPackage = append(headerPackage, bodyPacketCount...)

	buffer = append(buffer, headerPackage)

	ding1 := ""
	ding1 = task.GetResponseBody()
	if ding1 != "" {
		log.Info("")
	}

	// Body packets
	start := 0
	end := packetSize - 1
	for i := 0; i < packetCount; i++ {
		buffer = append(buffer, []byte(taskRespondBody[start:end]))
		start = end + 1
		end += packetSize
	}

	if lastPacket {
		start = len(taskRespondBody) - (len(taskRespondBody) - packetSize*packetCount)
		buffer = append(buffer, []byte(taskRespondBody[start:]))
	}

	return buffer
}
