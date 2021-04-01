package sherlockcrawler

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"errors"

	sherlockkafka "github.com/DerAlexx/SherlockGopher/sherlockkafka"
	"github.com/segmentio/kafka-go"
)

const (
	topictask         = "tasktopic"
	topicurl         = "urltopic"
	brokerAddress = "0.0.0.0:9092"
)

type KafkaWriter struct {
	writer kafka.Writer
}

func NewKafkaWriter(topic string, brokAddress string) *KafkaWriter {
	return &KafkaWriter{
		writer: kafka.Writer{
			Addr:  kafka.TCP(brokAddress),
			Topic: topic,
		},
	}
}

func convert(task *CrawlerTaskRequest) *sherlockkafka.KafkaTask {
	tmp := sherlockkafka.KafkaTask{}
	tmpmap := make(map[string][]string)
	for headerkey, headerValue := range *task.responseHeader {
		tmpmap[headerkey] = headerValue
	}
	tmp.TaskID =            task.taskID
	tmp.Addr = task.addr
	if (task.taskError == nil){
		tmperr := errors.New("")
		tmp.TaskError =   tmperr.Error()
	} else {
		tmp.TaskError =  task.taskError.Error()
	} 
	tmp.StatusCode =         task.statusCode
	tmp.ResponseBodyBytes =  task.responseBodyBytes
	tmp.ResponseHeader =     	tmpmap
	tmp.ResponseTime =       task.responseTime
	return &tmp
}

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
