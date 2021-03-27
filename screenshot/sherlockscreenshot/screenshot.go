package sherlockscreenshot

import (
	"context"
	"log"
	"encoding/json"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	sherlockkafka "github.com/DerAlexx/SherlockGopher/sherlockkafka"
	"github.com/segmentio/kafka-go"
)

const (
	topicConsume         = "urltoscreenshotservice"
	topicProduce         = "sendscreenshots"
	brokerAddress = "localhost:9092"
)

type Screenshot struct {
	Picture []byte
	URL     string
}

func StartChrome(url string) {
	_, cancel := chromedp.NewContext(context.Background())
	defer cancel()
}

func TakeScreenshot(url string) *Screenshot {

	url = "https://golangcode.com/"

	// List of actions to run in sequence (which also fills our image buffer)
	var imageBuf []byte
	if err := chromedp.Run(context.TODO(), ScreenshotTasks(url, &imageBuf)); err != nil {
		log.Fatal(err)
	}
	/* Write image to file
	filename := "golangcode.png"
	if err := ioutil.WriteFile(filename, imageBuf, 0644); err != nil {
		log.Fatal(err)
	}*/
	return &Screenshot{imageBuf, url}
}

/*
ScreenshotTasks creates a screenshot.
*/
func ScreenshotTasks(url string, imageBuf *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) (err error) {
			*imageBuf, err = page.CaptureScreenshot().WithQuality(90).Do(ctx)
			return err
		}),
	}
}

func (db *DB) ConsumeUrlForScreenshot(ctx context.Context) {
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic:   topicConsume,
	})
	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}
		// after receiving the message create task

		tmpurl := sherlockkafka.KafkaUrl{}
		err = json.Unmarshal(msg.Value, &tmpurl)
		if err != nil {
			panic("parsing json failed" + err.Error())
		}
		res := TakeScreenshot(tmpurl.URL)
		db.Save(res)
	}
}

//TODO: python env dinge installieren, celerey docker compose, kafka consumer producer