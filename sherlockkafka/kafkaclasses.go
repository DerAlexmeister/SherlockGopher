package sherlockkafka

import (
	"time"

)

type KafkaUrl struct {
	URL string `json:"url"`
}

type KafkaTask struct {
	TaskID            uint64         `json:"taskid"` 
	Addr              string         `json:"addr"`   
	TaskError         string          `json:"taskerror"`   
	ResponseHeader    map[string][]string    `json:"responseheader"` 
	ResponseBodyBytes []byte         `json:"responsebodybytes"`
	StatusCode        int            `json:"statuscode"`      
	ResponseTime      time.Duration  `json:"responstime"`      
}

type KafkaScreenshot struct {
	Picture []byte	`json:"picture"`
	URL     string	`json:"url"`
}

/*
KafkaRequestedURL will be the struct for the post request of the domain to crawl.
*/
type KafkaRequestedURL struct {
	URL string `json:"url" binding:"required"`
}

/*
KafkaRequestedURL will be the struct for the post request of the domain to crawl.
*/
type KafkaAkkRequestedURL struct {
	Status bool `json:"status" binding:"required"`
}

/*
KafkaRequestedStatus will be the struct for the post request of the status functions.
*/
type KafkaRequestedStatus struct {
	Operation string `json:"operation" binding:"required"`
	Target    string `json:"target" binding:"required"`
}


func NewKafkaUrl() *KafkaUrl {
	return &KafkaUrl{}
}

func NewKafkaTask() *KafkaTask {
	return &KafkaTask{}
}

func NewKafkaScreenshot() *KafkaScreenshot {
	return &KafkaScreenshot{}
}

/*
KafkaRequestedURL will be a new instance of KafkaRequestedURL.
*/
func NewKafkaRequestedURL() *KafkaRequestedURL {
	return &KafkaRequestedURL{}
}

/*
KafkaAkkRequestedURL will be a new instance of KafkaAkkRequestedURL.
*/
func NewAkkKafkaRequestedURL() *KafkaAkkRequestedURL {
	return &KafkaAkkRequestedURL{}
}

/*
KafkaRequestedStatus will be a new instance of KafkaRequestedStatus.
*/
func NewKafkaRequestedStatus() *KafkaRequestedStatus {
	return &KafkaRequestedStatus{}
}