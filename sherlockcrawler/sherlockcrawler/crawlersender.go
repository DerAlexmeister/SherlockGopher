package sherlockcrawler

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"

	sherlockkafka "github.com/DerAlexx/SherlockGopher/sherlockkafka"
	"github.com/segmentio/kafka-go"
)

const (
	topic         = "crawlertasks"
	brokerAddress = "localhost:9092"
)

/*
CrawlerTaskRequest will be a request made by the analyser.
*/
type KafkaWriter struct {
	writer kafka.Writer
}

func NewKafkaWriter() *KafkaWriter {
	return &KafkaWriter{
		writer: kafka.Writer{
			Addr:  kafka.TCP(brokerAddress),
			Topic: topic,
		},
	}
}

func (sherlock *SherlockCrawler) produce(ctx context.Context, task *CrawlerTaskRequest, wg *sync.WaitGroup) error {

	if task.GetTaskError() != nil {
		wg.Done()
	} else {
		(*task).setTaskState(PROCESSING)

		tmp := &sherlockkafka.KafkaTask{
			TaskID:            task.GetTaskID(),
			Addr:              task.GetAddr(),
			TaskState:         task.GetTaskState(),
			TaskError:         task.GetTaskError(),
			TaskErrorTry:      task.GetTryError(),
			Response:          task.response,
			ResponseHeader:    task.responseHeader,
			ResponseBody:      task.GetResponseBody(),
			ResponseBodyBytes: task.GetResponseBodyInBytes(),
			StatusCode:        task.GetStatusCode(),
			ResponseTime:      task.GetResponseTime()}

		res1B, _ := json.Marshal(tmp)

		// each kafka message has a key and value. The key is used
		// to decide which partition (and consequently, which broker)
		// the message gets published on
		err := sherlock.kwriter.writer.WriteMessages(ctx, kafka.Message{
			Key: []byte(strconv.Itoa(0)),
			// create an arbitrary message payload for the value
			Value: res1B,
		})
		if err != nil {
			panic("could not write message " + err.Error())
		}
		wg.Done()
	}
	return nil
}

func consume(ctx context.Context) {
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
	})
	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}
		// after receiving the message create task
		var stringurl string
		url := json.Unmarshal(msg.Value, &stringurl)
	}
}
