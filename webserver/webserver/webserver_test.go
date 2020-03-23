package webserver

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

func TestSendHelloPing(t *testing.T) { // TODO drüber schauen.

	req, err := http.NewRequest("GET", "localhost:8080/helloping", nil)
	if err != nil {
		log.Fatal(err)

	}
	log.Printf("Body %s", req.Body)
	t.Log(req)
}

func TestSubmitURL(t *testing.T) {
	req, err := http.NewRequest("POST", "/search", nil)
	if err != nil {
		log.Fatal(err)

	}
	fmt.Println(req)
}
