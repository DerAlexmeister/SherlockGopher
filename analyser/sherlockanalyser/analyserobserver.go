package sherlockanalyser

import (
	log "github.com/sirupsen/logrus"
)

//SystemState stores the current state of the analyser service (STOP, PAUSE, RUNNING, IDLE).
type SystemState int

const (
	//STOP System is shut down. Kill and remove all queue tasks. No further actions.
	STOP = 0

	//PAUSE System is paused. Pause task execution.
	PAUSE = 1

	//RUNNING System is running. Start all tasks.
	RUNNING = 2

	//IDLE System is idle. Waiting for tasks with low performance.
	IDLE = 4
)

/*
AnalyserObserverInterface is the interface of the AnalyserMediator.
*/
type AnalyserObserverInterface interface {
	UpdateState(state int)
	CleanCommand()
}

/*
CleanCommand cleans the analyser queue.
*/
func (obs AnalyserObserver) CleanCommand() {
	log.Info("Observer called CleanQueue")
	obs.queue.CleanQueue()
}

/*
UpdateState setter for the state of the service/queue.
*/
func (obs AnalyserObserver) UpdateState(state int) {
	log.Info("Observer updated state")
	if !(state > 4 || state < 0) {
		if obs.Analyser() != nil {
			obs.analyser.SetState(state)
		}
		if obs.Queue() != nil {
			obs.queue.SetState(state)
		}
	}
}

/*
AnalyserObserver mediates important values of the analyser system.
*/
type AnalyserObserver struct {
	analyser *AnalyserServiceHandler
	queue    *AnalyserQueue
}

/*
Queue returns the queue of the analyser.
*/
func (obs *AnalyserObserver) Queue() *AnalyserQueue {
	return obs.queue
}

/*
SetQueue is a setter for the queue of the analyser.
*/
func (obs *AnalyserObserver) SetQueue(queue *AnalyserQueue) {
	obs.queue = queue
	queue.SetObserver(obs)
}

/*
Analyser returns the analyser stored in the observer.
*/
func (obs *AnalyserObserver) Analyser() *AnalyserServiceHandler {
	return obs.analyser

}

/*
SetAnalyser is a setter for the analyser stored in the observer.
*/
func (obs *AnalyserObserver) SetAnalyser(analyser *AnalyserServiceHandler) {
	obs.analyser = analyser
	analyser.SetObserver(obs)
}

/*
NewAnalyserObserver creates a new analyser oberver and returns it.
*/
func NewAnalyserObserver() AnalyserObserver {
	return AnalyserObserver{}
}
