package main

import (
	"context"

	screenshot "github.com/DerAlexx/SherlockGopher/screenshot/sherlockscreenshot"
)

func main() {
	scrser := screenshot.NewScreenshotService()
	db := screenshot.Connect()
	db.ConsumeUrlForScreenshot(context.TODO(), scrser.GetContext())
}
