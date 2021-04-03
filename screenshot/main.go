package main

import (
	"context"

	screenshot "github.com/DerAlexx/SherlockGopher/screenshot/sherlockscreenshot"
)

func main() {
	scrser := screenshot.NewScreenshotService()
	scrser.ConsumeUrlForScreenshot(context.TODO())
}
