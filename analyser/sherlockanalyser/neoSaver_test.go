package sherlockanalyser

import (
	"fmt"
	"net/http"
	"testing"

	neo "github.com/DerAlexx/SherlockGopher/sherlockneo"
	"github.com/DerAlexx/SherlockGopher/sherlockneo/mocks"
	"github.com/golang/mock/gomock"
)

func TestSave(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	session := mocks.NewMockSession(mockCtrl)
	driver := mocks.NewMockDriver(mockCtrl)
	result := mocks.NewMockResult(mockCtrl)
	record := mocks.NewMockRecord(mockCtrl)

	driver.EXPECT().Session(gomock.Any()).Return(session, nil).MinTimes(1)
	driver.EXPECT().Close().MinTimes(1)
	session.EXPECT().Run(gomock.Any(), nil).Return(result, nil).MinTimes(1)
	session.EXPECT().Close().Return(nil).MinTimes(1)
	result.EXPECT().Record().Return(record).MinTimes(1)
	result.EXPECT().Err().Return(fmt.Errorf("test")).MinTimes(1)
	result.EXPECT().Next().Return(true).MinTimes(1)
	var ret interface{} = false
	record.EXPECT().Get(gomock.Any()).Return(ret, true).MinTimes(1)

	saver := newNeoSaver()
	saver.session = session
	saver.driver = driver

	header := http.Header{}
	cd := CrawlerData{
		responseHeader: &header,
		statusCode:     200,
		responseTime:   200,
		taskError:      fmt.Errorf("test"),
	}

	task := NewTask(&cd)
	task.saver = &saver
	task.saver.SetDriver(driver)
	links := make([]string, 2)
	links[0] = "test.com"
	links[1] = "."
	task.SetFoundLinks(links)
	task.SetWorkAddr("workAddr.com")
	task.fileType = neo.HTML

	saver.Save(task)
}

func TestGetSession(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	session := mocks.NewMockSession(mockCtrl)
	saver := newNeoSaver()
	saver.session = session

	if saver.GetSession() == nil {
		t.Fatal("GetSession failed")
	}
}

func TestContains(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	session := mocks.NewMockSession(mockCtrl)
	driver := mocks.NewMockDriver(mockCtrl)
	result := mocks.NewMockResult(mockCtrl)
	record := mocks.NewMockRecord(mockCtrl)

	session.EXPECT().Run(gomock.Any(), nil).Return(result, nil)
	result.EXPECT().Record().Return(record)
	result.EXPECT().Next().Return(true)
	var ret interface{} = false
	record.EXPECT().Get(gomock.Any()).Return(ret, true)

	saver := newNeoSaver()
	saver.session = session
	saver.driver = driver

	task := NewTask(nil)
	links := make([]string, 1)
	links[0] = "test.com"
	task.SetFoundLinks(links)
	if len(saver.Contains(task)) != 1 {
		t.Fatal("Contains failed")
	}
}
