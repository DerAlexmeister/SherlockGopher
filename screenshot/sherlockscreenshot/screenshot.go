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
	topicUrl        = "urltopic"
	brokerAddress = "0.0.0.0:9092"
)

type Screenshot struct {
	Picture []byte
	URL     string
}

type ScreenshotService struct {
	Chromecontext context.Context
}

func NewScreenshot() *Screenshot {
	return &Screenshot{}
}

func NewScreenshotService() *ScreenshotService {
	ctx := startChrome()
	screenservice := ScreenshotService{
		Chromecontext:        ctx,
	}
	return &screenservice
}

func startChrome() context.Context{
	ctx, cancel := chromedp.NewContext(context.TODO())
	defer cancel()
	return ctx
}


func (scr *Screenshot) setPicture(pic []byte) {
	scr.Picture = pic
}

func (scr *Screenshot) setUrl(url string) {
	scr.URL = url
}

func (scr *Screenshot) getPicture() *[]byte {
	return &scr.Picture
}

func (scr *Screenshot) getUrl() string {
	return scr.URL
}

func (scrser *ScreenshotService) GetContext() context.Context {
	return scrser.Chromecontext
}

func TakeScreenshot(url string, chromectx context.Context) *Screenshot {

	var imageBuf []byte
	if err := chromedp.Run(chromectx, ScreenshotTasks(url, &imageBuf)); err != nil {
		log.Fatal(err)
	}

	tmpscr := NewScreenshot()
	tmpscr.setPicture(imageBuf)
	tmpscr.setUrl(url)
	return tmpscr
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

func (db *DB) ConsumeUrlForScreenshot(ctx context.Context, chromectx context.Context) {
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic:   topicUrl,
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
		res := TakeScreenshot(tmpurl.URL, chromectx)
		db.Save(res)
	}
}