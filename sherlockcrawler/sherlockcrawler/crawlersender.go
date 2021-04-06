package sherlockcrawler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"

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
convert creates a new KafkaTask with a CrawlerTaskRequest
*/
func convert(task *CrawlerTaskRequest) *sherlockkafka.KafkaTask {
	tmp := sherlockkafka.KafkaTask{}
	tmpmap := make(map[string][]string)
	for headerkey, headerValue := range *task.responseHeader {
		tmpmap[headerkey] = headerValue
	}
	tmp.TaskID = task.taskID
	tmp.Addr = task.addr
	if task.taskError == nil {
		tmperr := errors.New("")
		tmp.TaskError = tmperr.Error()
	} else {
		tmp.TaskError = task.taskError.Error()
	}
	tmp.StatusCode = task.statusCode
	tmp.ResponseBodyBytes = task.responseBodyBytes
	tmp.ResponseHeader = tmpmap
	tmp.ResponseTime = task.responseTime
	return &tmp
}

/*
SendTaskToAnalyser is a kafka producer, sending tasks to the analyser
*/
func (sherlock *SherlockCrawler) SendTaskToAnalyser(ctx context.Context, task *CrawlerTaskRequest, wg *sync.WaitGroup) error {

	if task.GetTaskError() != nil {
		wg.Done()
	} else {
		(*task).setTaskState(PROCESSING)

		tmp := convert(task)

		res1B, _ := json.Marshal(tmp)
		fmt.Println(res1B)

		kwriter := NewKafkaWriter(topictask, brokerAddress)

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
		wg.Done()
	}
	return nil
}

/*
ReceiveUrlFromAnalyser is a kafka consumer, receiving urls from the crawler
*/
func (sherlock *SherlockCrawler) ReceiveUrlFromAnalyser(ctx context.Context) {
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic:   topicurl,
	})
	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}
		// after receiving the message create task
		//stringurl := sherlockkafka.KafkaUrl{}
		var tmpurl sherlockkafka.KafkaUrl
		err = json.Unmarshal(msg.Value, &tmpurl)
		if err != nil {
			panic("parsing json failed" + err.Error())
		}
		sherlock.NextCreateTask(tmpurl.URL)
	}
}
