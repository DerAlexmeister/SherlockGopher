# Docs about the analyser.

The analyser service serves 3 purposes. It contains the receiving part of the streaming service,
the parser and its functionality and the counterpart of the crawler which handles tasks.
The parser and the streaming service are described in a seperate doc.

This doc describes the handling of tasks in the analyser.
After the streaming service receives a chunk it creates a new task and appends it to the analyser queue.

A task in the analser service looks like this:

```go
type AnalyserTaskRequest struct {
	id           uint64
	workAddr     string
	domainInfo   *jp.URL
	html         string
	state        TaskState
	linkTags     map[string]string
	fileType     string
	foundLinks   []string
	parserTime   int64
	analyserTime int64
	crawlerData  *CrawlerData
	Dependencies *AnalyserDependency
}
```

A received chunk is stored in the field crawlerData. The id is used to identify a task and make it unique.
The task status is changed dependent on the current state of processing.

States of the AnalyserTaskRequest:
	
* UNDONE: a task has not been started yet.	
* PROCESSING: a task has been started.	
* CRAWLERERROR: a task contains a crawler error.	
* SAVING: a task is saving data to NEO4J.	
* SENDTOCRAWLER: a task is sending data to the crawler.
* FINISHED: a task is finished and ready to be remove from the queue.

The current number of tasks in a specific state are counted in the QueueStatus. This statistic can be displayed in the frontend.

```go
type QueueStatus struct {
	undoneTasks        uint64
	processingTasks    uint64
	crawlerErrorTasks  uint64
	savingTasks        uint64
	sendToCrawlerTasks uint64
	finishedTasks      uint64
}
```

The AnalyserQueue has a status which is identical to the status of the analyser service.
States of the AnalyserQueue:
*   Stop: The analyser service will shut down. Kills and removes all queue tasks. No further actions will be executed.
*   Pause: The analyser service is still running but it doesnt accept new tasks. So it will complete its remaining tasks.
*   Clean: removes all tasks from the queue
*   Idle: if the queue is empty for a specific time the service switches to the idle state. In idle state the time before the crawlerqueue gets checked is higher. This decreases the performance usage.

The queue and and the service state is observed and updated by a observer. This functionality is implemented identically in the crawler service.
It works as described in the observer pattern.





