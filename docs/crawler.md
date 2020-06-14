# Docs about the crawler.

## Whats the SherlockCrawlerService

The main task of sherlock-crawler is to request a website, wait for the response, extract the needed information and stream it to the analyser.

At the start the webserver will submit the start address via gRPC to the crawler. From there on the crawler will make a HTTP-Request as well as starting a timer. The task of the timer will be, to measure the time between the request and the response. This is considered as RTT and will be stored in the responseTime. The Response will be stored in a struct of the type  ``` CrawlerTaskRequest ```.

```go
    type CrawlerTaskRequest struct {
        taskID            uint64 
        addr              string 
        taskState         taskState
        taskError         error 
        taskErrorTry      int   
        response          *http.Response
        responseHeader    *http.Header 
        responseBody      string
        responseBodyBytes []byte        
        statusCode        int           
        responseTime      time.Duration 
    }
```

* The task id is used to identify a ``` CrawlerTaskRequest ``` and understand which CrawlerTaskRequests belong together
* The adress stores the address of the website that was crawled.
* The taskState stores the current state of the task.
* The taskError stores errors that were caused by the http package.
* The taskErrorTry contains the number of errors that occurred.
* The responseHeader stores the header of the response.
* The responseBody stores the body of the response as a string.
* responseBodyBytes stores the body of the response as a byte array.
* statusCode represents the status code of the request.

A task knows four different states:
* undone: the taks wasnt processed yet
* processing: the taks is currently processed
* finished: the task was processed succesfully
* failed: an error occurred during processing

CrawlerTaskRequests will be stored in a queue, so they can be processed one after another. 
The CrawlerQueue has a status which is identical to the status of the crawler service.
The streaming service will take tasks from the queue and stream them to the analyser.

The crawler service can be set in four different states: 
*   Stop: The crawler will shut down. Kills and removes all queue tasks. No further actions will be executed.
*   Pause: The crawler is still running but it doesnt accept new tasks. So it will complete its remaining tasks.
*   Clean: removes all tasks from the queue
*   Idle: if the queue is empty for a specific time the service switches to the idle state. In idle state the time before the crawlerqueue gets checked is higher. This decreases the performance usage.


The queue and and the service state is observed and updated by a observer. This functionality is implemented identically in the analyser service.
It works as described in the observer pattern.