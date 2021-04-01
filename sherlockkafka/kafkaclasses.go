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