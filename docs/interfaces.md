# Description of all interfaces used

## Communication between Crawler and Analyser

The Crawler Service crawls websites and saves the most important information in a CrawlerTaskRequest

```go
type CrawlerTaskRequest struct {
	taskiD            uint64 //important
	addr              string //important
	taskstate         TASKSTATE
	taskerror         error //important
	taskerrortry      int   
	response          *http.Response
	responseHeader    *http.Header //important
	responseBody      string
	responseBodyBytes []byte   //important     
	statuscode        int //important          
	responseTime      time.Duration //important
}
```

* Every information highlighted with //important needs to be send to the Analyser for further processing.

* The transmission is implemented with grpc streams. The reason for that is the field responseBodyBytes.

* Dependant on the website this byte array could be to big for regular grpc transmission.

* We use a fixed chunksize to slice the byte array in small chunks

```go
# Proto File
message Chunk {
        bytes Content = 1;
        uint64 TaskId = 2;
        string Address = 3;
        string TaskError = 4;
        repeated HeaderArray Header = 5;
        int32 StatusCode = 6;
        int64 ResponseTime = 7;
}
```

* We use the same message for 3 different cases:
    * Something went wrong with a http response - in this case the field task error contains an error
    * Tranmission of the important information - the field task error is nil
    * Transmission of the byte array via small chunks - the field task error is nil

* In case there is no error:
    * First we send all important information except the byte array to the analyser
    * Afterwards we send the small chunks to the analyser where we put it back together
    * If everything went right the analyser will respond  with a status message "ok" (ack)

* In case there is an error:
    * The Crawler sends only 4 field to the Analyser: TaskId, Address, TaskError, ResponseTime

* The analyser checks the field  TaskError to distinguish which fiels of the message it needs to extract. 

* The extracted information will be saved in struct called analyserTaskRequest which is similar to the CrawlerTaskRequest
    

## Communication between Webserver and Crawler/Analyser

* The frontend gives the user the oportunity to interact with the services
* There are the following functions:
    * stop the analyser, stop the crawler or stop both services at once
    * pause the analyser, pause the crawler or pause both services at once
    * resume the analyser, resume the crawler or resume both services at once after a pause
    * clean the analyser, clean the crawler or clean both services at once. This will empty the queue filled tasks of the service.

* The frontend sends a post request to the webservice. The webservice will analyze the post request. Depending on the result the webservice sends a grpc message (StateRequest) to the service, where the functions are implemented. As soon as the service recieves the message it will respons with a messsage (State Response).

* Underneath you can see the implementation of the grpc messages in the proto file

```go
message StateRequest {
    CurrentState State = 1;
}

message StateResponse {
    bool received = 1;
}

enum CurrentState {
        Stop = 0;
        Pause = 1;
        Resume = 2;
        Clean = 3;
}
```