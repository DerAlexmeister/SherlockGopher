package sherlockscreenshot

import (
	"context"
	"log"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type Screenshot struct {
	Picture []byte
	URL     string
}

func StartChrome(url string) context.Context {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	return ctx
}

func TakeScreenshot(url string, ctx context.Context) *Screenshot {

	url = "https://golangcode.com/"

	// List of actions to run in sequence (which also fills our image buffer)
	var imageBuf []byte
	if err := chromedp.Run(ctx, ScreenshotTasks(url, &imageBuf)); err != nil {
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
