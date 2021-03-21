package sherlockscreenshot

import (
	"context"
	"log"
	"sync"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

/*
ScreenshotTaskRequest will be a request made by the screenshot service.
*/
type ScreenshotTaskRequest struct {
	id    uint64
	url   string
	state TaskState
}

/*
TaskState will be a type representing the current TaskState of the task.
*/
type TaskState int

const (
	//UNDONE A task has not been started yet.
	UNDONE TaskState = 0

	//PROCESSING A task has been started.
	PROCESSING TaskState = 1

	//FINISHED A task is finished and ready to be remove from the queue.
	FINISHED TaskState = 3
)

/*
getID returns Id.
*/
func (screenshotTask *ScreenshotTaskRequest) getID() uint64 {
	return screenshotTask.id
}

/*
SetID sets Id.
*/
func (screenshotTask *ScreenshotTaskRequest) SetID(id uint64) {
	screenshotTask.id = id
}

/*
getURL returns Id.
*/
func (screenshotTask *ScreenshotTaskRequest) getURL() string {
	return screenshotTask.url
}

/*
SetURL sets Id.
*/
func (screenshotTask *ScreenshotTaskRequest) SetURL(url string) {
	screenshotTask.url = url
}

/*
State returns State.
*/
func (screenshotTask *ScreenshotTaskRequest) getState() TaskState {
	return screenshotTask.state
}

/*
SetState sets State.
*/
func (screenshotTask *ScreenshotTaskRequest) SetState(state TaskState) {
	screenshotTask.state = state
}

func takeScreenshot(url string) *Screenshot {

	// Start Chrome
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	url = "https://golangcode.com/"

	// Run Tasks
	// List of actions to run in sequence (which also fills our image buffer)
	var imageBuf []byte
	if err := chromedp.Run(ctx, ScreenshotTasks(url, &imageBuf)); err != nil {
		log.Fatal(err)
	}
	/* Write our image to file
	filename := "golangcode.png"
	if err := ioutil.WriteFile(filename, imageBuf, 0644); err != nil {
		log.Fatal(err)
	}*/
	return &Screenshot{imageBuf, url}
}

/*
ScreenshotTasks creates a screenshot.
*/
func ScreenshotTasks(url string, imageBuf *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) (err error) {
			*imageBuf, err = page.CaptureScreenshot().WithQuality(90).Do(ctx)
			return err
		}),
	}
}

/*
Execute will execute the task.
*/
func (screenshotTask *ScreenshotTaskRequest) Execute(waitGroup *sync.WaitGroup) bool {
	res := takeScreenshot(screenshotTask.getURL())

	screenshotTask.SetState(FINISHED)

	if waitGroup != nil {
		defer waitGroup.Done()
	}

	return true
}

/*
NewTask returns a new task initialized by crawler data.
*/
func NewTask() *ScreenshotTaskRequest {
	task := ScreenshotTaskRequest{}
	return &task
}
