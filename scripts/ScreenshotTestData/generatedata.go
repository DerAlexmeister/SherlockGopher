package main

import (
	"fmt"

	screenshot "github.com/DerAlexx/SherlockGopher/screenshot/sherlockscreenshot"
)

func main() {
	service := screenshot.NewScreenshotService()
	/*for i := 0; i < 2; i++ {
		tmp := service.TakeScreenshot("https://golangcode.com/")
		service.GetClient().Save(tmp)
	}*/
	allscreenshots, err := service.GetClient().ReturnAllScreenshots()
	fmt.Println(allscreenshots, err)
	fmt.Println(len(allscreenshots))
}
