package main

import(
	screenshot "github.com/DerAlexx/SherlockGopher/screenshot/sherlockscreenshot"
	"context"
)

func main() {
	db := screenshot.Connect()
	go db.ConsumeUrlForScreenshot(context.TODO())
}
