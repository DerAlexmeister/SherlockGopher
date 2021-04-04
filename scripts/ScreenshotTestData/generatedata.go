package main

import (
	"fmt"

	screenshot "github.com/DerAlexx/SherlockGopher/screenshot/sherlockscreenshot"
)

func main() {
	service := screenshot.NewScreenshotService()
	defer service.GetCancelContext()
	for i := 0; i < 2; i++ {
		tmp := service.TakeScreenshot("https://golangcode.com/")
		fmt.Println(tmp)
		//service.GetClient().Save(tmp)
	}
	/*
		allscreenshots, err := service.GetClient().ReturnAllScreenshots()
		fmt.Println(allscreenshots, err)
		fmt.Println(len(allscreenshots))
	*/
}
