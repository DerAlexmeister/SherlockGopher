package main

import(
	screenshot "github.com/DerAlexx/SherlockGopher/screenshot/sherlockscreenshot"
)

func main() {
	db := screenshot.Connect()
	ctx := context.Background()
	db.client.ConsumeUrlForScreenshot(ctx)
}
