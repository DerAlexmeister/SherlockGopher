package sherlockkafka

import (
	"net/http"
	"time"
)

/*
AnalyserTaskResponse will be a response to the crawler with the url
*/
type KafkaUrl struct {
	URL string `json:"url"`
}

type KafkaTask struct {
	TaskID            uint64                    `json:"taskid"` //taskID, send every time.
	Addr              string                    `json:"addr"`   //addr, once
	TaskState         sherlockcrawler.taskState `json:"taskstate"`
	TaskError         error                     `json:"taskerror"`   //error, send as string in case there is an error then dont send a body
	TaskErrorTry      int                       `json:"taskerrotry"` //never
	Response          *http.Response            `json:"response"`
	ResponseHeader    *http.Header              `json:"responseheader"` //header, once (typ map)
	ResponseBody      string                    `json:"responsebody"`
	ResponseBodyBytes []byte                    `json:"responsebodybytes"` //body, split
	StatusCode        int                       `json:"statuscode"`        //statusCode, once
	ResponseTime      time.Duration             `json:"responstime"`       //response time, once
}
