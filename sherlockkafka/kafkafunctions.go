package sherlockkafka

import (
	"errors"
	"fmt"

	sherlockcrawler "github.com/DerAlexx/SherlockGopher/sherlockcrawler"
	"github.com/asaskevich/govalidator"
	"github.com/micro/go-micro/util/log"
)

/*
CreateTask will append the current queue with a task.
*/
func (sherlock sherlockcrawler.SherlockCrawler) NextCreateTask(url string) error {
	message := fmt.Sprintf("malformed or invalid url: %s", url)
	if isValid := govalidator.IsURL(url); isValid {
		task := sherlockcrawler.NewTask()
		task.setAddr(url)
		if id := sherlock.getQueue().AppendQueue(&task); id > 0 {
			log.Debug("Created task ", task.GetTaskID(), task.GetAddr())
			return nil
		}
		log.Error("Could not append task ", task.GetTaskID(), task.GetAddr())
	}
	return errors.New(message)
}
