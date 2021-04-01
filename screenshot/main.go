package main

import(
	screenshot "github.com/DerAlexx/SherlockGopher/screenshot/sherlockscreenshot"
	"context"
)

func main() {
	scrser := screenshot.NewScreenshotService()
	db := screenshot.Connect()
	go db.ConsumeUrlForScreenshot(context.TODO(), scrser.GetContext())
}
