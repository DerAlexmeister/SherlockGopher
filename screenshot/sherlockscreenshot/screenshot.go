package sherlockscreenshot

import (
	"context"
	"encoding/json"

	sherlockkafka "github.com/DerAlexx/SherlockGopher/sherlockkafka"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/segmentio/kafka-go"
)

const (
	topicUrl      = "urltopic"
	brokerAddress = "0.0.0.0:9092"
)

type Screenshot struct {
	Picture []byte
	URL     string
}

type ScreenshotService struct {
	Chromecontext       context.Context
	Chromecancelcontext context.CancelFunc
	Client              *DB
}

func NewScreenshot() *Screenshot {
	return &Screenshot{}
}

func NewScreenshotService() *ScreenshotService {
	ctx, ctxcancel := startChrome()
	client := Connect()

	screenservice := ScreenshotService{
		Chromecontext:       ctx,
		Chromecancelcontext: ctxcancel,
		Client:              client,
	}
	return &screenservice
}

func startChrome() (context.Context, context.CancelFunc) {
	/*opts := []chromedp.ExecAllocatorOption{
		chromedp.ExecPath("../../chromium/chromedriver"),
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.TODO(), opts...)

	ctx, _ := chromedp.NewContext(allocCtx)*/

	ctx, cancel := chromedp.NewContext(context.TODO())
	return ctx, cancel
}

func (scr *Screenshot) setPicture(pic []byte) {
	scr.Picture = pic
}

func (scr *Screenshot) setUrl(url string) {
	scr.URL = url
}

func (scr *Screenshot) GetPicture() *[]byte {
	return &scr.Picture
}

func (scr *Screenshot) GetUrl() string {
	return scr.URL
}

func (scrser *ScreenshotService) GetContext() context.Context {
	return scrser.Chromecontext
}

func (scrser *ScreenshotService) GetCancelContext() context.CancelFunc {
	return scrser.Chromecancelcontext
}

func (scrser *ScreenshotService) GetClient() *DB {
	return scrser.Client
}

func (scrser *ScreenshotService) TakeScreenshot(url string) *Screenshot {

	var imageBuf []byte
	if err := chromedp.Run(scrser.GetContext(), ScreenshotTasks(url, &imageBuf)); err != nil {
		panic(err)
	}
	defer scrser.GetCancelContext()

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
			*imageBuf, err = page.CaptureScreenshot().WithQuality(60).Do(ctx)
			return err
		}),
	}
}

func (scrser *ScreenshotService) ConsumeUrlForScreenshot(ctx context.Context) {
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
		res := scrser.TakeScreenshot(tmpurl.URL)
		scrser.GetClient().Save(res)
	}
}
