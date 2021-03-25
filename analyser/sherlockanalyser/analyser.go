package sherlockanalyser

import (
	"context"
	"errors"
	"sync"
	"time"

	swd "github.com/DerAlexx/SherlockGopher/sherlockwatchdog"
	log "github.com/sirupsen/logrus"

	proto "github.com/DerAlexx/SherlockGopher/analyser/proto"
)

const (
	startDelay = 200000
)

/*
AnalyserServiceHandler will be AnalyserService representation.
*/
//nolint: structcheck
type AnalyserServiceHandler struct {
	AnalyserQueue *AnalyserQueue
	Dependencies  *AnalyserDependency
	observer      AnalyserObserverInterface
	watchdog      swd.WatchdogInterface
	state         int
	delay         *time.Duration
	idleCounter   int
	kwriter       *KafkaWriter
}

/*
Watchdog returns a analyser watchdog.
*/
func (analyser *AnalyserServiceHandler) Watchdog() swd.WatchdogInterface {
	return analyser.watchdog
}

/*
SetObserver is a setter for the analyser oberver.
*/
func (analyser *AnalyserServiceHandler) SetObserver(observer AnalyserObserverInterface) {
	analyser.observer = observer
}

/*
Observer returns a analyser observer.
*/
func (analyser *AnalyserServiceHandler) Observer() AnalyserObserverInterface {
	return analyser.observer
}

/*
State returns the state of the analyser.
*/
func (analyser *AnalyserServiceHandler) State() int {
	return analyser.state
}

/*
SetState sets the state of the analyser.
*/
func (analyser *AnalyserServiceHandler) SetState(state int) {
	analyser.state = state
}

/*
Delay returns the delay time.
*/
func (analyser *AnalyserServiceHandler) Delay() *time.Duration {
	return analyser.delay
}

/*
SetDelay sets the delay time.
*/
func (analyser *AnalyserServiceHandler) SetDelay(delay *time.Duration) {
	analyser.delay = delay
}

/*
InjectDependency will inject the dependencies for the a analyser instance
into the actual instance.
*/
func (analyser *AnalyserServiceHandler) InjectDependency(deps *AnalyserDependency) {
	analyser.Dependencies = deps
}

/*
NewAnalyserServiceHandler will return an new AnalyserServiceHandler instance.
*/
func NewAnalyserServiceHandler() *AnalyserServiceHandler {
	que := NewAnalyserQueue()
	watchdog := swd.NewSherlockWatchdog()
	que.SetWatchdog(&watchdog)
	observer := NewAnalyserObserver()
	kwriter := NewKafkaWriter()
	analyser := AnalyserServiceHandler{
		AnalyserQueue: &que,
		Dependencies:  nil,
		observer:      observer,
		watchdog:      &watchdog,
		kwriter:       kwriter,
	}

	observer.SetAnalyser(&analyser)
	observer.SetQueue(&que)

	observer.UpdateState(RUNNING)
	return &analyser
}

/*
getQueue will return a pointer to the queue.
*/
func (analyser *AnalyserServiceHandler) getQueue() *AnalyserQueue {
	return analyser.AnalyserQueue
}

/*
manageUndoneTasks will manage all undone tasks and start them.
*/
func (analyser *AnalyserServiceHandler) manageUndoneTasks() {
	var localWaitGroup sync.WaitGroup

	for _, v := range *(*analyser.getQueue()).getThisQueue() {
		if v.State() == UNDONE && analyser.State() == RUNNING {
			v.InjectDependency(analyser.Dependencies)
			v.SetState(PROCESSING)
			localWaitGroup.Add(1)
			go v.Execute(&localWaitGroup, analyser)
		}
	}
	localWaitGroup.Wait()
}

/*
manageFinishedTasks will managed all finished tasks.
*/
func (analyser *AnalyserServiceHandler) manageFinishedTasks() {
	for k, v := range *(*analyser.getQueue()).getThisQueue() {
		if v.State() == FINISHED {
			(*analyser.getQueue()).RemoveFromQueue(k)
		}
	}
}

/*
runManager will run all tasks related functions on the analyserServiceHandler and put them into a go routine.
*/
func (analyser *AnalyserServiceHandler) runManager() {
	analyser.manageUndoneTasks()
	analyser.manageFinishedTasks()
}

/*
ManageTasks will be a function to check all tasks in a period of time.
Should run as a go routine. Will work like a cronJob.
Attention: The for and the if are important!
*/
func (analyser *AnalyserServiceHandler) ManageTasks() {
	log.Debug("Started ManageTasks")

	for analyser.State() != STOP {
		if analyser.State() != PAUSE {
			analyser.runManager()
		}
		analyser.Watchdog().Watch()
	}
}

/*
WorkloadRPC can be called to get the workload of the analyser
*/
func (analyser *AnalyserServiceHandler) WorkloadRPC(ctx context.Context, _ *proto.WorkloadRequest, out *proto.WorkloadResponse) error {
	log.Debug("Called WorkloadRPC")
	undone, processing, crawlerError, saving, sendToCrawler, finished := analyser.getQueue().getNumberOfStatus()
	out.Undone = undone
	out.Processing = processing
	out.CrawlerError = crawlerError
	out.Saving = saving
	out.SendToCrawler = sendToCrawler
	out.Finished = finished
	return nil
}

/*
StateRPC will get the status of the analyser.
*/
func (analyser *AnalyserServiceHandler) StateRPC(ctx context.Context, in *proto.StateRequest, out *proto.StateResponse) error {
	log.Debug("Called StateRPC")
	out.State = proto.AnalyserStateEnum(analyser.State())
	return nil
}

/*
ChangeStateRPC can be called to set the state of the analyser
The idle state is not allowed to be set manually.
*/
func (analyser *AnalyserServiceHandler) ChangeStateRPC(ctx context.Context, in *proto.ChangeStateRequest, out *proto.ChangeStateResponse) error {
	var err error = nil
	out.Success = true

	switch in.State {
	case proto.AnalyserStateEnum_Stop:
		analyser.Observer().UpdateState(STOP)
	case proto.AnalyserStateEnum_Pause:
		analyser.Observer().UpdateState(PAUSE)
	case proto.AnalyserStateEnum_Running:
		analyser.Observer().UpdateState(RUNNING)
	case proto.AnalyserStateEnum_Clean:
		analyser.Observer().CleanCommand()
	case proto.AnalyserStateEnum_Idle:
		out.Success = false
		errMessage := "idle state is not allowed to be set manually"
		err = errors.New(errMessage)
		log.Warn(errMessage)
	default:
		out.Success = false
		errMessage := "unexpected State"
		err = errors.New(errMessage)
		log.Warn(err)
	}

	return err
}
