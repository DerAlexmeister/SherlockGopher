package sherlockanalyser

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	sherlockkafka "github.com/DerAlexx/SherlockGopher/sherlockkafka"
	"github.com/segmentio/kafka-go"
)

// kafka topics and broker address
var brokerAddress, topictask, topicurl string

/*
KafkaWriter to write into kafak topics
*/
type KafkaWriter struct {
	writer kafka.Writer
}

/*
NewKafkaWriter creates a new KafakWriter instance
*/
func NewKafkaWriter(topic string, brokAddress string) *KafkaWriter {
	return &KafkaWriter{
		writer: kafka.Writer{
			Addr:  kafka.TCP(brokAddress),
			Topic: topic,
		},
	}
}

/*
Init prepares urls for kafka
*/
func Init() {
	brokerAddress = readFromENV("KAFKA_BROKER", "0.0.0.0:9092")
	topictask = readFromENV("KAFKA_TOPIC_TASK", "testtask")
	topicurl = readFromENV("KAFKA_TOPIC_URL", "testurl")
}

/*
readFromENV allows docker usage
*/
func readFromENV(key, defaultVal string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultVal
	}
	return value
}

/*
convertUrlToKafkaUrl creates a new KafkaUrl with a url
*/
func convertUrlToKafkaUrl(url string) *sherlockkafka.KafkaUrl {
	tmp := sherlockkafka.KafkaUrl{}
	tmp.URL = url
	return &tmp
}

/*
SendUrlToCrawler is a kafka producer, sending urls to the crawler
*/
func (sherlock *AnalyserServiceHandler) SendUrlToCrawler(ctx context.Context, url string) error {

	kTask := convertUrlToKafkaUrl(url)
	res1B, _ := json.Marshal(&kTask)
	fmt.Println(res1B)
	kwriter := NewKafkaWriter(topicurl, brokerAddress)

	// each kafka message has a key and value. The key is used
	// to decide which partition (and consequently, which broker)
	// the message gets published on
	err := kwriter.writer.WriteMessages(ctx, kafka.Message{
		Key: []byte(strconv.Itoa(0)),
		// create an arbitrary message payload for the value
		Value: res1B,
	})
	if err != nil {
		panic("could not write message " + err.Error())
	}
	return nil
}

/*
SendUrlToCrawler is a kafka consumer, receiving tasks from the crawler
*/
func (analyser *AnalyserServiceHandler) ReceiveTaskFromCrawler(ctx context.Context) {
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic:   topictask,
	})
	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}
		// after receiving the message create task
		fmt.Println(msg)

		var tmptask sherlockkafka.KafkaTask
		err = json.Unmarshal(msg.Value, &tmptask)
		if err != nil {
			panic("parsing json failed" + err.Error())
		}

		headerMap := http.Header{}
		for headerKey, headerValue := range tmptask.ResponseHeader {
			for _, headVal := range headerValue {
				headerMap.Add(headerKey, headVal)
			}
		}

		task := NewCrawlerData()
		task.setTaskID(tmptask.TaskID)
		task.setAddr(tmptask.Addr)
		task.setTaskError(errors.New(tmptask.TaskError))
		task.setResponseHeader(&headerMap)
		task.setResponseBody(tmptask.ResponseBodyBytes)
		task.setStatusCode(tmptask.StatusCode)
		task.setResponseTime(tmptask.ResponseTime)
		analyserTaskRequest := NewTask(task)

		analyser.getQueue().AppendQueue(analyserTaskRequest)
	}
}
