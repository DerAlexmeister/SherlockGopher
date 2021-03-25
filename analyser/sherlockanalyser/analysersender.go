package sherlockanalyser

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	sherlockkafka "github.com/DerAlexx/SherlockGopher/sherlockkafka"
	"github.com/segmentio/kafka-go"
)

const (
	topic         = "test1"
	brokerAddress = "localhost:9092"
)

/*
KafkaWriter.
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

func (sherlock *AnalyserServiceHandler) produce(ctx context.Context, url string) error {

	tmp := &sherlockkafka.KafkaUrl{
		URL: url,
	}

	res1B, _ := json.Marshal(tmp)
	fmt.Println(res1B)

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
	return nil
}

func (analyser *AnalyserServiceHandler) Consume(ctx context.Context) {
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
		fmt.Println(msg)

		tmptask := sherlockkafka.KafkaTask{}
		err = json.Unmarshal(msg.Value, &tmptask)
		if err != nil {
			panic("parsing json failed" + err.Error())
		}

		task := NewCrawlerData()
		task.setTaskID(tmptask.TaskID)
		task.setAddr(tmptask.Addr)
		task.setTaskError(tmptask.TaskError)
		task.setResponseHeader(tmptask.ResponseHeader)
		task.setResponseBody(tmptask.ResponseBodyBytes)
		task.setStatusCode(tmptask.StatusCode)
		task.setResponseTime(tmptask.ResponseTime)
		analyserTaskRequest := NewTask(task)

		analyser.getQueue().AppendQueue(analyserTaskRequest)
	}
}
