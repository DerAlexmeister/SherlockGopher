package sherlockcrawler

import (
	"context"
	"fmt"
	"sync"
	"time"

	swd "github.com/ob-algdatii-20ss/SherlockGopher/sherlockwatchdog"
	log "github.com/sirupsen/logrus"

	"github.com/asaskevich/govalidator"
	aproto "github.com/ob-algdatii-20ss/SherlockGopher/analyser/proto"
	proto "github.com/ob-algdatii-20ss/SherlockGopher/sherlockcrawler/proto"
	"github.com/pkg/errors"
)

//Time to wait in milliseconds after checking for the tasks again.
const (
	maxErrors  = 3
	startDelay = 20000
)

/*
SherlockCrawler will be the crawlerService.
*/
//nolint: structcheck
type SherlockCrawler struct {
	StreamQueue     CrawlerQueueInterface //Queue with all tasks
	Dependencies    *SherlockDependencies
	Queue           *CrawlerQueue //Queue with all tasks
	observer        CrawlerObserverInterface
	watchdog        swd.WatchdogInterface
	analyserService *aproto.AnalyserService
	state           int
	delay           *time.Duration
	idleCounter     int
	kwriter         *KafkaWriter
}

/*
Watchdog will return the Watchdog.
Watchdog returns the watchdog of the SherlockCrawler.
*/
func (sherlock *SherlockCrawler) Watchdog() swd.WatchdogInterface {
	return sherlock.watchdog
}

/*
SetObserver will set the opbserver of a CrawlerObserver.
SetObserver is a setter for the watchdog of the SherlockCrawler.
*/
func (sherlock *SherlockCrawler) SetObserver(observer CrawlerObserverInterface) {
	sherlock.observer = observer
}

/*
Observer will return the Observer.
Observer returns the observer of the SherlockCrawler.
*/
func (sherlock *SherlockCrawler) Observer() CrawlerObserverInterface {
	return sherlock.observer
}

/*
State returns the state of the sherlock.
*/
func (sherlock *SherlockCrawler) State() int {
	return sherlock.state
}

/*
SetCrawlerState sets the state of the sherlock.
*/
func (sherlock *SherlockCrawler) SetCrawlerState(state int) {
	sherlock.state = state
}

/*
Delay returns the delay time.
*/
func (sherlock *SherlockCrawler) Delay() *time.Duration {
	return sherlock.delay
}

/*
SetDelay sets the delay time.
*/
func (sherlock *SherlockCrawler) SetDelay(delay *time.Duration) {
	sherlock.delay = delay
}

/*
SherlockDependencies is an type to manage all dependencies of sherlockCrawler.
*/
type SherlockDependencies struct {
	Analyser func() aproto.AnalyserService
}

/*
NewSherlockDependencies will return a new sherlockDependencies instance to put it in the dependencies
in a sherlockCrawler object.
*/
func NewSherlockDependencies() *SherlockDependencies {
	return &SherlockDependencies{}
}

/*
NewSherlockCrawlerService will return a new sherlockCrawler instance.
*/
func NewSherlockCrawlerService() *SherlockCrawler {
	que := NewCrawlerQueue()
	watchdog := swd.NewSherlockWatchdog()
	que.SetWatchdog(&watchdog)
	observer := NewCrawlerObserver()
	streamQue := NewCrawlerQueue()
	kwriter := NewKafkaWriter()

	crawler := SherlockCrawler{
		Queue:        &que,
		Dependencies: nil,
		observer:     observer,
		StreamQueue:  &streamQue,
		watchdog:     &watchdog,
		kwriter:      kwriter,
	}

	observer.SetCrawler(&crawler)
	observer.SetQueue(&que)
	observer.SetStreamerQue(&streamQue)

	observer.UpdateState(RUNNING)

	return &crawler
}

/*
InjectDependency will inject the dependencies for the a sherlockCrawler instance
into the actual instance.
*/
func (sherlock *SherlockCrawler) InjectDependency(deps *SherlockDependencies) {
	sherlock.Dependencies = deps
	serv := sherlock.Dependencies.Analyser()
	sherlock.analyserService = &serv
}

/*
getDependency will return a pointer to the dependencies instance of this service.
*/
func (sherlock *SherlockCrawler) getDependency() *SherlockDependencies {
	return sherlock.Dependencies
}

/*
getQueue will return a pointer to the queue.
*/
func (sherlock *SherlockCrawler) getQueue() *CrawlerQueue {
	return sherlock.Queue
}

/*
getStreamQueue will return a pointer to the queue.
*/
func (sherlock *SherlockCrawler) getStreamQueue() CrawlerQueueInterface {
	return sherlock.StreamQueue
}

/*
CreateTask will append the current queue with a task.
*/
func (sherlock SherlockCrawler) CreateTask(ctx context.Context, in *proto.CrawlTaskCreateRequest, out *proto.CrawlTaskCreateResponse) error {
	message := fmt.Sprintf("malformed or invalid url: %s", in.GetUrl())
	if isValid := govalidator.IsURL(in.GetUrl()); isValid {
		task := NewTask()
		task.setAddr(in.GetUrl())
		if id := sherlock.getQueue().AppendQueue(&task); id > 0 {
			out.Statuscode = proto.URL_STATUS_ok
			out.Taskid = task.GetTaskID()
			log.Debug("Created task ", task.GetTaskID(), task.GetAddr())
			return nil
		}
		log.Error("Could not append task ", task.GetTaskID(), task.GetAddr())
	}
	out.Statuscode = proto.URL_STATUS_failure
	out.Taskid = 0
	return errors.New(message)
}

/*
manageUndoneTasks will manage all undone tasks and start them.
*/
func (sherlock *SherlockCrawler) manageUndoneTasks(wg *sync.WaitGroup) {
	var localWaitGroup sync.WaitGroup
	for _, v := range *sherlock.getQueue().GetThisQueue() {
		if v.GetTaskState() == UNDONE {
			localWaitGroup.Add(1)
			go v.MakeRequestAndStoreResponse(&localWaitGroup)
			localWaitGroup.Wait()
		}
	}
	wg.Done()
}

/*
manageFinishedTasks will managed all finished tasks.
*/
//nolint: errcheck
func (sherlock *SherlockCrawler) manageFinishedTasks(wg *sync.WaitGroup) {
	var sendWG sync.WaitGroup
	for k, v := range *sherlock.getQueue().GetThisQueue() {
		if v.GetTaskState() == FINISHED {
			sendWG.Add(1)
			go sherlock.SendWebsiteData(sherlock.Dependencies.Analyser(), v, &sendWG)
			sendWG.Wait()
			sherlock.getQueue().RemoveFromQueue(k)
		}
	}
	wg.Done()
}

/*
manageFailedTasks manage all tasks which failed.
*/
//nolint: errcheck
func (sherlock *SherlockCrawler) manageFailedTasks(wg *sync.WaitGroup) {
	var localWaitGroup sync.WaitGroup
	var sendWG sync.WaitGroup
	for k, v := range *sherlock.getQueue().GetThisQueue() {
		if v.GetTaskState() == FAILED {
			log.WithFields(log.Fields{
				"addr": v.GetAddr(),
			}).Info("manageFailedTasks")
			if v.GetTryError() < maxErrors {
				localWaitGroup.Add(1)
				go v.MakeRequestAndStoreResponse(&localWaitGroup)
				localWaitGroup.Wait()
			} else {
				sendWG.Add(1)
				go sherlock.SendWebsiteData(sherlock.Dependencies.Analyser(), v, &sendWG)
				sendWG.Wait()
				sherlock.getQueue().RemoveFromQueue(k)
			}
		}
	}
	wg.Done()
}

func (sherlock *SherlockCrawler) runManager() {
	log.Info("Started RunManager")
	var wg sync.WaitGroup
	wg.Add(1)
	go sherlock.manageUndoneTasks(&wg)
	wg.Wait()
	wg.Add(1)
	go sherlock.manageFinishedTasks(&wg)
	wg.Wait()
	wg.Add(1)
	go sherlock.manageFailedTasks(&wg)
	wg.Wait()
}

/*
ManageTasks will be a function to check all tasks in a period of time.
Should run as a go routine. Will work like a cronJob.
*/
//nolint:staticcheck
func (sherlock *SherlockCrawler) ManageTasks() {
	var sleepTime time.Duration = 2
	log.Info("Started ManageTasks")
	for sherlock.State() != STOP {
		if sherlock.State() != PAUSE {
			sherlock.runManager()
		}
		time.Sleep(sleepTime)
		sherlock.Watchdog().Watch()
	}
}

/*
ReceiveURL will spawn the first task in the queue in order to start the howl process.
*/
func (sherlock *SherlockCrawler) ReceiveURL(ctx context.Context, in *proto.SubmitURLRequest, out *proto.SubmitURLResponse) error {
	var err = errors.New("malformed URL, please submit a well-formed one")
	if isValid := govalidator.IsURL(in.GetURL()); isValid {
		log.Info("Received URL from web server, will result in a new task with address .", in.GetURL())
		task := NewTask()
		task.setAddr(in.GetURL())
		sherlock.getQueue().AppendQueue(&task)
		out.Recieved = true
		return nil
	}

	log.Error("Received URL from web server, will not result in a new task because given address is not valid ", in.GetURL())
	out.Recieved = false
	out.Error = err.Error()
	return err
}

/*
StatusOfTaskQueue will send the status of the queue.
*/
func (sherlock *SherlockCrawler) StatusOfTaskQueue(ctx context.Context, _ *proto.TaskStatusRequest, out *proto.TaskStatusResponse) error {
	undone, processing, finished, failed := sherlock.getQueue().getNumberOfStatus()
	out.Undone = undone
	out.Processing = processing
	out.Finished = finished
	out.Failed = failed
	return nil
}

/*
GetState will get the status of the sherlock.
*/
func (sherlock *SherlockCrawler) GetState(_ context.Context, _ *proto.StateGetRequest, out *proto.StateGetResponse) error {
	out.State = proto.CurrentState(sherlock.CrawlerState())
	return nil
}

/*
SetState will set the status of the sherlock.
*/
func (sherlock *SherlockCrawler) SetState(_ context.Context, in *proto.StateRequest, out *proto.StateResponse) error {
	switch in.State {
	case proto.CurrentState_Stop:
		sherlock.Observer().UpdateState(STOP)
	case proto.CurrentState_Pause:
		sherlock.Observer().UpdateState(PAUSE)
	case proto.CurrentState_Running:
		sherlock.Observer().UpdateState(RUNNING)
	case proto.CurrentState_Clean:
		sherlock.getQueue().CleanQueue()
		sherlock.getStreamQueue().CleanQueue()
	default:
		out.Received = false
		log.Error("Crawler->SetState->Unknown state ", in.State)
		return errors.New("Unknown state")
	}
	out.Received = true

	return nil
}

/*
CrawlerState will return the current state of the crawler
*/
func (sherlock *SherlockCrawler) CrawlerState() int {
	return sherlock.state
}
