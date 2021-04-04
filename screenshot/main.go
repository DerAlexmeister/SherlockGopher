package main

import (
	"context"

	screenshot "github.com/DerAlexx/SherlockGopher/screenshot/sherlockscreenshot"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Started screenshot")
	screenshot.Init()
	scrser := screenshot.NewScreenshotService()
	scrser.ConsumeUrlForScreenshot(context.TODO())
}
