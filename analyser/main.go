package main

import (
	"fmt"
	"github.com/ob-algdatii-20ss/leistungsnachweis-dievierausrufezeichen/analyser/html2treeparser/model"
)

func main() {
	fmt.Println("test")
	rootNode := model.NewHTMLTree("<html></html>").Parse()
}
