package sherlockkafka

import (
	"time"
)

/*
KafkaUrl is used for the kafka service to send url between the screenshot, crawler and analyser service
*/
type KafkaUrl struct {
	URL string `json:"url"`
}

/*
KafkaUrl is used for the kafka service to send tasks between the crawler and analyser service
*/
type KafkaTask struct {
	TaskID            uint64              `json:"taskid"`
	Addr              string              `json:"addr"`
	TaskError         string              `json:"taskerror"`
	ResponseHeader    map[string][]string `json:"responseheader"`
	ResponseBodyBytes []byte              `json:"responsebodybytes"`
	StatusCode        int                 `json:"statuscode"`
	ResponseTime      time.Duration       `json:"responstime"`
}
