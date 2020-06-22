package sherlockwatchdog

import (
	"fmt"
	"sync"
	"time"
)

/*
WatchdogInterface which contains all functions for the watchdog.
*/
type WatchdogInterface interface {
	Delay() time.Duration
	Fass()
	Aus()
	Watch()
}

const (
	maxWaitTime = 2000 * time.Millisecond
	minWaitTime = 2 * time.Microsecond
	hitsToCatch = 100
)

/*
SherlockWatchdog is the struct representing the watchdog.
*/
type SherlockWatchdog struct {
	delay *time.Duration
	mutex *sync.Mutex
	hits  int
}

/*
NewSherlockWatchdog creates a new SherlockWatchdog and returns it.
*/
func NewSherlockWatchdog() SherlockWatchdog {
	delay := minWaitTime
	mutex := sync.Mutex{}
	return SherlockWatchdog{
		delay: &delay,
		mutex: &mutex,
		hits:  0,
	}
}

/*
Delay returns the delay.
*/
func (swd *SherlockWatchdog) Delay() time.Duration {
	return *swd.delay
}

/*
Fass increments the hits (+1).
*/
func (swd *SherlockWatchdog) Fass() {
	swd.mutex.Lock()
	swd.hit()
	swd.mutex.Unlock()
}

/*
Aus sets the hits to zero.
*/
func (swd *SherlockWatchdog) Aus() {
	swd.mutex.Lock()
	swd.reset()
	swd.mutex.Unlock()
}

/*
Watch checks if the watchdog needs to be slown down.
*/
func (swd *SherlockWatchdog) Watch() {
	swd.mutex.Lock()
	swd.waitAndHit()
	swd.mutex.Unlock()
}

func (swd *SherlockWatchdog) hit() {
	swd.hits++
}

func (swd *SherlockWatchdog) reset() {
	swd.hits = 0
	*swd.delay = minWaitTime
	fmt.Println("Working...")
}

func (swd *SherlockWatchdog) waitAndHit() {
	time.Sleep(swd.Delay())
	swd.hit()
	if swd.hits >= hitsToCatch {
		swd.slow()
	}
}

func (swd *SherlockWatchdog) slow() {
	if *swd.delay < maxWaitTime {
		*swd.delay *= 10
	} else {
		fmt.Println("Sleeping...")
	}
}
