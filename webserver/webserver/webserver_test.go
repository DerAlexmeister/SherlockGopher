package webserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	analyserproto "github.com/DerAlexx/SherlockGopher/analyser/proto"
	crawlerproto "github.com/DerAlexx/SherlockGopher/sherlockcrawler/proto"
	"github.com/DerAlexx/SherlockGopher/sherlockneo/mocks"
	mockanalyser "github.com/DerAlexx/SherlockGopher/webserver/webserver/mocks/analyser"
	mockcrawler "github.com/DerAlexx/SherlockGopher/webserver/webserver/mocks/crawler"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

const (
	amountofnodes = "amountofnodes"
)

type WebserverTestError struct{}

func (err *WebserverTestError) Error() string {
	return "This error is expected!"
}

func TestSendHelloPing(t *testing.T) { // TODO drüber schauen.
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	sut := New()
	sut.Helloping(c)

	assert.Equal(t, 200, w.Result().StatusCode)
}

//nolint: misspell
func TestReceiveUrl(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/search", bytes.NewReader([]byte("{\n        \"url\": \"www.github.com/\"\n    } ")))
	sut := New()
	ctrl := gomock.NewController(t)
	crawler := mockcrawler.NewMockCrawlerService(ctrl)
	crawler.EXPECT().ReceiveURL(gomock.Any(), gomock.Any()).Return(&crawlerproto.SubmitURLResponse{Recieved: true}, nil).Times(1)
	analyser := mockanalyser.NewMockAnalyserService(ctrl)
	sut.Dependency = &SherlockWebServerDependency{
		Crawler: func() crawlerproto.CrawlerService {
			return crawler
		},
		Analyser: func() analyserproto.AnalyserService {
			return analyser
		},
	}
	sut.ReceiveURL(c)
	fmt.Println(w.Body)
	var got gin.H
	err := json.Unmarshal(w.Body.Bytes(), &got)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "200 OK", w.Result().Status)
	assert.Equal(t, "Fine", got["Status"])
}

func TestReceiveMetadata(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/meta", nil)
	sut := New()
	ctrl := gomock.NewController(t)
	crawler := mockcrawler.NewMockCrawlerService(ctrl)
	analyser := mockanalyser.NewMockAnalyserService(ctrl)
	analyser.EXPECT().WorkloadRPC(gomock.Any(), gomock.Any()).Return(&analyserproto.WorkloadResponse{
		CrawledWebsite: "",
		Undone:         0,
		Processing:     0,
		CrawlerError:   0,
		Saving:         0,
		SendToCrawler:  0,
		Finished:       0,
	}, nil).Times(1)
	crawler.EXPECT().StatusOfTaskQueue(gomock.Any(), gomock.Any()).Return(&crawlerproto.TaskStatusResponse{
		Website:    "",
		Undone:     0,
		Processing: 0,
		Finished:   0,
		Failed:     0,
	}, nil).Times(1)
	sut.Dependency = &SherlockWebServerDependency{
		Crawler: func() crawlerproto.CrawlerService {
			return crawler
		},
		Analyser: func() analyserproto.AnalyserService {
			return analyser
		},
	}

	sut.ReceiveMetadata(c)

	var got gin.H
	err := json.Unmarshal(w.Body.Bytes(), &got)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "200 OK", w.Result().Status)
	assert.Equal(t, "map[Analyser:map[CrawlerError:0 Finished:0 Processing:0 Saving:0 SendToCrawler:0 Undone:0 Website:] Crawler:map[Failed:0 Finished:0 Processing:0 Undone:0 Website:]]", fmt.Sprint(got))
}

func TestReceiveMetadataCrawlerAnalyserError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/meta", nil)
	sut := New()
	ctrl := gomock.NewController(t)
	crawler := mockcrawler.NewMockCrawlerService(ctrl)
	analyser := mockanalyser.NewMockAnalyserService(ctrl)
	analyser.EXPECT().WorkloadRPC(gomock.Any(), gomock.Any()).Return(&analyserproto.WorkloadResponse{
		CrawledWebsite: "",
		Undone:         0,
		Processing:     0,
		CrawlerError:   0,
		Saving:         0,
		SendToCrawler:  0,
		Finished:       0,
	}, &WebserverTestError{}).Times(1)
	crawler.EXPECT().StatusOfTaskQueue(gomock.Any(), gomock.Any()).Return(&crawlerproto.TaskStatusResponse{
		Website:    "",
		Undone:     0,
		Processing: 0,
		Finished:   0,
		Failed:     0,
	}, &WebserverTestError{}).Times(1)
	sut.Dependency = &SherlockWebServerDependency{
		Crawler: func() crawlerproto.CrawlerService {
			return crawler
		},
		Analyser: func() analyserproto.AnalyserService {
			return analyser
		},
	}

	sut.ReceiveMetadata(c)

	var got gin.H
	err := json.Unmarshal(w.Body.Bytes(), &got)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "500 Internal Server Error", w.Result().Status)
	assert.Equal(t, "Couldnt get Metadata, analyser and crawler services are unavailable", fmt.Sprint(got["Status"]))
}

func TestReceiveMetadataAnalyserError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/meta", nil)
	sut := New()
	ctrl := gomock.NewController(t)
	crawler := mockcrawler.NewMockCrawlerService(ctrl)
	analyser := mockanalyser.NewMockAnalyserService(ctrl)
	analyser.EXPECT().WorkloadRPC(gomock.Any(), gomock.Any()).Return(&analyserproto.WorkloadResponse{
		CrawledWebsite: "",
		Undone:         0,
		Processing:     0,
		CrawlerError:   0,
		Saving:         0,
		SendToCrawler:  0,
		Finished:       0,
	}, &WebserverTestError{}).Times(1)
	crawler.EXPECT().StatusOfTaskQueue(gomock.Any(), gomock.Any()).Return(&crawlerproto.TaskStatusResponse{
		Website:    "",
		Undone:     0,
		Processing: 0,
		Finished:   0,
		Failed:     0,
	}, nil).Times(1)
	sut.Dependency = &SherlockWebServerDependency{
		Crawler: func() crawlerproto.CrawlerService {
			return crawler
		},
		Analyser: func() analyserproto.AnalyserService {
			return analyser
		},
	}

	sut.ReceiveMetadata(c)

	var got gin.H
	err := json.Unmarshal(w.Body.Bytes(), &got)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "500 Internal Server Error", w.Result().Status)
	assert.Equal(t, "Couldnt get Metadata, analyser service is unavailable", fmt.Sprint(got["Status"]))
}

func TestReceiveMetadataCrawlerError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/meta", nil)
	sut := New()
	ctrl := gomock.NewController(t)
	crawler := mockcrawler.NewMockCrawlerService(ctrl)
	analyser := mockanalyser.NewMockAnalyserService(ctrl)
	analyser.EXPECT().WorkloadRPC(gomock.Any(), gomock.Any()).Return(&analyserproto.WorkloadResponse{
		CrawledWebsite: "",
		Undone:         0,
		Processing:     0,
		CrawlerError:   0,
		Saving:         0,
		SendToCrawler:  0,
		Finished:       0,
	}, nil).Times(1)
	crawler.EXPECT().StatusOfTaskQueue(gomock.Any(), gomock.Any()).Return(&crawlerproto.TaskStatusResponse{
		Website:    "",
		Undone:     0,
		Processing: 0,
		Finished:   0,
		Failed:     0,
	}, &WebserverTestError{}).Times(1)
	sut.Dependency = &SherlockWebServerDependency{
		Crawler: func() crawlerproto.CrawlerService {
			return crawler
		},
		Analyser: func() analyserproto.AnalyserService {
			return analyser
		},
	}

	sut.ReceiveMetadata(c)

	var got gin.H
	err := json.Unmarshal(w.Body.Bytes(), &got)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "500 Internal Server Error", w.Result().Status)
	assert.Equal(t, "Couldnt get Metadata, crawler service is unavailable", fmt.Sprint(got["Status"]))
}

//nolint: misspell
func TestReceiveUrlBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	sut := New()
	ctrl := gomock.NewController(t)
	crawler := mockcrawler.NewMockCrawlerService(ctrl)
	crawler.EXPECT().ReceiveURL(gomock.Any(), gomock.Any()).Return(&crawlerproto.SubmitURLResponse{Recieved: true}, nil).Times(1)
	analyser := mockanalyser.NewMockAnalyserService(ctrl)
	sut.Dependency = &SherlockWebServerDependency{
		Crawler: func() crawlerproto.CrawlerService {
			return crawler
		},
		Analyser: func() analyserproto.AnalyserService {
			return analyser
		},
	}
	sut.ReceiveURL(c)
	assert.Equal(t, 400, w.Result().StatusCode)
}

//nolint: misspell
func TestReceiveUrlReceivedFalse(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/search", bytes.NewReader([]byte("{\n        \"url\": \"www.github.com/\"\n    }")))
	sut := New()
	ctrl := gomock.NewController(t)
	crawler := mockcrawler.NewMockCrawlerService(ctrl)
	crawler.EXPECT().ReceiveURL(gomock.Any(), gomock.Any()).Return(&crawlerproto.SubmitURLResponse{Recieved: false}, nil).Times(1)
	analyser := mockanalyser.NewMockAnalyserService(ctrl)
	sut.Dependency = &SherlockWebServerDependency{
		Crawler: func() crawlerproto.CrawlerService {
			return crawler
		},
		Analyser: func() analyserproto.AnalyserService {
			return analyser
		},
	}
	sut.ReceiveURL(c)

	var got gin.H
	err := json.Unmarshal(w.Body.Bytes(), &got)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "200 OK", w.Result().Status)
	assert.Equal(t, "The webserver cannot submit your URL to the Crawler, couldnt reach the crawler service", got["Message"])
}

func TestChangeStateCrawler(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/meta", bytes.NewReader([]byte("{\n        \"operation\": \"Clean\",\n        \"target\": \"Crawler\"\n    }")))
	sut := New()
	ctrl := gomock.NewController(t)
	crawler := mockcrawler.NewMockCrawlerService(ctrl)

	crawler.EXPECT().StatusOfTaskQueue(gomock.Any(), gomock.Any()).Return(&crawlerproto.TaskStatusResponse{
		Website:    "",
		Undone:     0,
		Processing: 0,
		Finished:   0,
		Failed:     0,
	}, nil).Times(1)
	sut.Dependency = &SherlockWebServerDependency{
		Crawler: func() crawlerproto.CrawlerService {
			return crawler
		},
	}
	crawler.EXPECT().SetState(gomock.Any(), gomock.Any()).Return(nil, nil)

	sut.ChangeState(c)

	var got gin.H
	err := json.Unmarshal(w.Body.Bytes(), &got)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "200 OK", w.Result().Status)
	assert.Equal(t, "Fine", fmt.Sprint(got["Status"]))
}

func TestChangeStateAnalyser(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/meta", bytes.NewReader([]byte("{\n        \"operation\": \"Resume\",\n        \"target\": \"Analyser\"\n    }")))
	sut := New()
	ctrl := gomock.NewController(t)

	analyser := mockanalyser.NewMockAnalyserService(ctrl)
	analyser.EXPECT().WorkloadRPC(gomock.Any(), gomock.Any()).Return(&analyserproto.WorkloadResponse{
		CrawledWebsite: "",
		Undone:         0,
		Processing:     0,
		CrawlerError:   0,
		Saving:         0,
		SendToCrawler:  0,
		Finished:       0,
	}, nil).Times(1)

	sut.Dependency = &SherlockWebServerDependency{
		Analyser: func() analyserproto.AnalyserService {
			return analyser
		},
	}

	analyser.EXPECT().ChangeStateRPC(gomock.Any(), gomock.Any()).Return(nil, nil)

	sut.ChangeState(c)

	var got gin.H
	err := json.Unmarshal(w.Body.Bytes(), &got)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "200 OK", w.Result().Status)
	assert.Equal(t, "Fine", fmt.Sprint(got["Status"]))
}

//nolint: dupl
func TestChangeStateAnalyserAndCrawler(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/meta", bytes.NewReader([]byte("{\n        \"operation\": \"Resume\",\n        \"target\": \"All\"\n    }")))
	sut := New()
	ctrl := gomock.NewController(t)

	crawler := mockcrawler.NewMockCrawlerService(ctrl)
	analyser := mockanalyser.NewMockAnalyserService(ctrl)
	analyser.EXPECT().WorkloadRPC(gomock.Any(), gomock.Any()).Return(&analyserproto.WorkloadResponse{
		CrawledWebsite: "",
		Undone:         0,
		Processing:     0,
		CrawlerError:   0,
		Saving:         0,
		SendToCrawler:  0,
		Finished:       0,
	}, nil).Times(1)
	crawler.EXPECT().StatusOfTaskQueue(gomock.Any(), gomock.Any()).Return(&crawlerproto.TaskStatusResponse{
		Website:    "",
		Undone:     0,
		Processing: 0,
		Finished:   0,
		Failed:     0,
	}, nil).Times(1)
	sut.Dependency = &SherlockWebServerDependency{
		Crawler: func() crawlerproto.CrawlerService {
			return crawler
		},
		Analyser: func() analyserproto.AnalyserService {
			return analyser
		},
	}

	analyser.EXPECT().ChangeStateRPC(gomock.Any(), gomock.Any()).Return(nil, nil)
	crawler.EXPECT().SetState(gomock.Any(), gomock.Any()).Return(nil, nil)

	sut.ChangeState(c)

	var got gin.H
	err := json.Unmarshal(w.Body.Bytes(), &got)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "200 OK", w.Result().Status)
	assert.Equal(t, "Fine", fmt.Sprint(got["Status"]))
}

//nolint: dupl
func TestChangeStateAnalyserAndCrawlerStatusStop(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/meta", bytes.NewReader([]byte("{\n        \"operation\": \"Stop\",\n        \"target\": \"All\"\n    }")))
	sut := New()
	ctrl := gomock.NewController(t)

	crawler := mockcrawler.NewMockCrawlerService(ctrl)
	analyser := mockanalyser.NewMockAnalyserService(ctrl)
	analyser.EXPECT().WorkloadRPC(gomock.Any(), gomock.Any()).Return(&analyserproto.WorkloadResponse{
		CrawledWebsite: "",
		Undone:         0,
		Processing:     0,
		CrawlerError:   0,
		Saving:         0,
		SendToCrawler:  0,
		Finished:       0,
	}, nil).Times(1)
	crawler.EXPECT().StatusOfTaskQueue(gomock.Any(), gomock.Any()).Return(&crawlerproto.TaskStatusResponse{
		Website:    "",
		Undone:     0,
		Processing: 0,
		Finished:   0,
		Failed:     0,
	}, nil).Times(1)
	sut.Dependency = &SherlockWebServerDependency{
		Crawler: func() crawlerproto.CrawlerService {
			return crawler
		},
		Analyser: func() analyserproto.AnalyserService {
			return analyser
		},
	}

	analyser.EXPECT().ChangeStateRPC(gomock.Any(), gomock.Any()).Return(nil, nil)
	crawler.EXPECT().SetState(gomock.Any(), gomock.Any()).Return(nil, nil)

	sut.ChangeState(c)

	var got gin.H
	err := json.Unmarshal(w.Body.Bytes(), &got)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "200 OK", w.Result().Status)
	assert.Equal(t, "Fine", fmt.Sprint(got["Status"]))
}

//nolint: dupl
func TestChangeStateAnalyserAndCrawlerStatusPause(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/meta", bytes.NewReader([]byte("{\n        \"operation\": \"Pause\",\n        \"target\": \"All\"\n    }")))
	sut := New()
	ctrl := gomock.NewController(t)

	crawler := mockcrawler.NewMockCrawlerService(ctrl)
	analyser := mockanalyser.NewMockAnalyserService(ctrl)
	analyser.EXPECT().WorkloadRPC(gomock.Any(), gomock.Any()).Return(&analyserproto.WorkloadResponse{
		CrawledWebsite: "",
		Undone:         0,
		Processing:     0,
		CrawlerError:   0,
		Saving:         0,
		SendToCrawler:  0,
		Finished:       0,
	}, nil).Times(1)
	crawler.EXPECT().StatusOfTaskQueue(gomock.Any(), gomock.Any()).Return(&crawlerproto.TaskStatusResponse{
		Website:    "",
		Undone:     0,
		Processing: 0,
		Finished:   0,
		Failed:     0,
	}, nil).Times(1)
	sut.Dependency = &SherlockWebServerDependency{
		Crawler: func() crawlerproto.CrawlerService {
			return crawler
		},
		Analyser: func() analyserproto.AnalyserService {
			return analyser
		},
	}

	analyser.EXPECT().ChangeStateRPC(gomock.Any(), gomock.Any()).Return(nil, nil)
	crawler.EXPECT().SetState(gomock.Any(), gomock.Any()).Return(nil, nil)

	sut.ChangeState(c)

	var got gin.H
	err := json.Unmarshal(w.Body.Bytes(), &got)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "200 OK", w.Result().Status)
	assert.Equal(t, "Fine", fmt.Sprint(got["Status"]))
}

func TestChangeStateCrawlerUnknown(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/meta", bytes.NewReader([]byte("{\n        \"operation\": \"Clean\",\n        \"target\": \"Crawler\"\n    }")))
	sut := New()
	ctrl := gomock.NewController(t)
	crawler := mockcrawler.NewMockCrawlerService(ctrl)

	crawler.EXPECT().StatusOfTaskQueue(gomock.Any(), gomock.Any()).Return(&crawlerproto.TaskStatusResponse{
		Website:    "",
		Undone:     0,
		Processing: 0,
		Finished:   0,
		Failed:     0,
	}, nil).Times(1)
	sut.Dependency = &SherlockWebServerDependency{
		Crawler: func() crawlerproto.CrawlerService {
			return crawler
		},
	}
	crawler.EXPECT().SetState(gomock.Any(), gomock.Any()).Return(nil, &WebserverTestError{})

	sut.ChangeState(c)

	var got gin.H
	err := json.Unmarshal(w.Body.Bytes(), &got)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "200 OK", w.Result().Status)
	assert.Equal(t, "Unknown", fmt.Sprint(got["Crawler"]))
}

func TestChangeStateAnalyserUnknown(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/meta", bytes.NewReader([]byte("{\n        \"operation\": \"Resume\",\n        \"target\": \"Analyser\"\n    }")))
	sut := New()
	ctrl := gomock.NewController(t)

	analyser := mockanalyser.NewMockAnalyserService(ctrl)
	analyser.EXPECT().WorkloadRPC(gomock.Any(), gomock.Any()).Return(&analyserproto.WorkloadResponse{
		CrawledWebsite: "",
		Undone:         0,
		Processing:     0,
		CrawlerError:   0,
		Saving:         0,
		SendToCrawler:  0,
		Finished:       0,
	}, nil).Times(1)

	sut.Dependency = &SherlockWebServerDependency{
		Analyser: func() analyserproto.AnalyserService {
			return analyser
		},
	}

	analyser.EXPECT().ChangeStateRPC(gomock.Any(), gomock.Any()).Return(nil, &WebserverTestError{})

	sut.ChangeState(c)

	var got gin.H
	err := json.Unmarshal(w.Body.Bytes(), &got)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "200 OK", w.Result().Status)
	assert.Equal(t, "Unknown", fmt.Sprint(got["Analyser"]))
}

//nolint: dupl
func TestChangeStateAnalyserAndCrawlerUnknown(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/meta", bytes.NewReader([]byte("{\n        \"operation\": \"Resume\",\n        \"target\": \"All\"\n    }")))
	sut := New()
	ctrl := gomock.NewController(t)

	crawler := mockcrawler.NewMockCrawlerService(ctrl)
	analyser := mockanalyser.NewMockAnalyserService(ctrl)
	analyser.EXPECT().WorkloadRPC(gomock.Any(), gomock.Any()).Return(&analyserproto.WorkloadResponse{
		CrawledWebsite: "",
		Undone:         0,
		Processing:     0,
		CrawlerError:   0,
		Saving:         0,
		SendToCrawler:  0,
		Finished:       0,
	}, nil).Times(1)
	crawler.EXPECT().StatusOfTaskQueue(gomock.Any(), gomock.Any()).Return(&crawlerproto.TaskStatusResponse{
		Website:    "",
		Undone:     0,
		Processing: 0,
		Finished:   0,
		Failed:     0,
	}, nil).Times(1)
	sut.Dependency = &SherlockWebServerDependency{
		Crawler: func() crawlerproto.CrawlerService {
			return crawler
		},
		Analyser: func() analyserproto.AnalyserService {
			return analyser
		},
	}

	analyser.EXPECT().ChangeStateRPC(gomock.Any(), gomock.Any()).Return(nil, &WebserverTestError{})
	crawler.EXPECT().SetState(gomock.Any(), gomock.Any()).Return(nil, &WebserverTestError{})

	sut.ChangeState(c)

	var got gin.H
	err := json.Unmarshal(w.Body.Bytes(), &got)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "200 OK", w.Result().Status)
	assert.Equal(t, "Unknown", fmt.Sprint(got["Analyser"]))
	assert.Equal(t, "Unknown", fmt.Sprint(got["Crawler"]))
}

func TestChangeStateBadJSON(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/meta", bytes.NewReader([]byte("{\n        \"operation\"")))
	sut := New()

	sut.ChangeState(c)

	assert.Equal(t, "400 Bad Request", w.Result().Status)
	assert.Equal(t, "{\"Status\":\"Error while reveiving Requested Status\"}{\"Status\":{}}{\"Status\":\"Unknown Target, expected Crawler, Analyser or All\"}", w.Body.String())
}

func TestGetServiceStatus(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/meta", bytes.NewReader([]byte("{\n        \"operation\": \"Running\",\n        \"target\": \"All\"\n    }")))
	sut := New()
	ctrl := gomock.NewController(t)

	crawler := mockcrawler.NewMockCrawlerService(ctrl)
	analyser := mockanalyser.NewMockAnalyserService(ctrl)
	analyser.EXPECT().WorkloadRPC(gomock.Any(), gomock.Any()).Return(&analyserproto.WorkloadResponse{
		CrawledWebsite: "",
		Undone:         0,
		Processing:     0,
		CrawlerError:   0,
		Saving:         0,
		SendToCrawler:  0,
		Finished:       0,
	}, nil).Times(1)
	crawler.EXPECT().StatusOfTaskQueue(gomock.Any(), gomock.Any()).Return(&crawlerproto.TaskStatusResponse{
		Website:    "",
		Undone:     0,
		Processing: 0,
		Finished:   0,
		Failed:     0,
	}, nil).Times(1)
	sut.Dependency = &SherlockWebServerDependency{
		Crawler: func() crawlerproto.CrawlerService {
			return crawler
		},
		Analyser: func() analyserproto.AnalyserService {
			return analyser
		},
	}

	analyser.EXPECT().StateRPC(gomock.Any(), gomock.Any()).Return(&analyserproto.StateResponse{
		State: analyserproto.AnalyserStateEnum_Clean,
	}, nil)
	crawler.EXPECT().GetState(gomock.Any(), gomock.Any()).Return(&crawlerproto.StateGetResponse{
		State: crawlerproto.CurrentState_Running,
	}, nil)

	sut.GetServiceStatus(c)

	var got gin.H
	err := json.Unmarshal(w.Body.Bytes(), &got)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "200 OK", w.Result().Status)
	assert.Equal(t, "Clean", fmt.Sprint(got["Analyser"]))
	assert.Equal(t, "Running", fmt.Sprint(got["Crawler"]))
}

//nolint: dupl
func TestGetServiceStatusAnalyserCrawlerError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/meta", bytes.NewReader([]byte("{\n        \"operation\": \"Resume\",\n        \"target\": \"All\"\n    }")))
	sut := New()
	ctrl := gomock.NewController(t)

	crawler := mockcrawler.NewMockCrawlerService(ctrl)
	analyser := mockanalyser.NewMockAnalyserService(ctrl)
	analyser.EXPECT().WorkloadRPC(gomock.Any(), gomock.Any()).Return(&analyserproto.WorkloadResponse{
		CrawledWebsite: "",
		Undone:         0,
		Processing:     0,
		CrawlerError:   0,
		Saving:         0,
		SendToCrawler:  0,
		Finished:       0,
	}, nil).Times(1)
	crawler.EXPECT().StatusOfTaskQueue(gomock.Any(), gomock.Any()).Return(&crawlerproto.TaskStatusResponse{
		Website:    "",
		Undone:     0,
		Processing: 0,
		Finished:   0,
		Failed:     0,
	}, nil).Times(1)
	sut.Dependency = &SherlockWebServerDependency{
		Crawler: func() crawlerproto.CrawlerService {
			return crawler
		},
		Analyser: func() analyserproto.AnalyserService {
			return analyser
		},
	}

	analyser.EXPECT().StateRPC(gomock.Any(), gomock.Any()).Return(nil, &WebserverTestError{})
	crawler.EXPECT().GetState(gomock.Any(), gomock.Any()).Return(nil, &WebserverTestError{})

	sut.GetServiceStatus(c)

	var got gin.H
	err := json.Unmarshal(w.Body.Bytes(), &got)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "200 OK", w.Result().Status)
	assert.Equal(t, "Unknown", fmt.Sprint(got["Analyser"]))
	assert.Equal(t, "Unknown", fmt.Sprint(got["Crawler"]))
}

func TestGetServiceStatusCrawlerError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/meta", bytes.NewReader([]byte("{\n        \"operation\": \"Resume\",\n        \"target\": \"All\"\n    }")))
	sut := New()
	ctrl := gomock.NewController(t)

	crawler := mockcrawler.NewMockCrawlerService(ctrl)
	analyser := mockanalyser.NewMockAnalyserService(ctrl)
	analyser.EXPECT().WorkloadRPC(gomock.Any(), gomock.Any()).Return(&analyserproto.WorkloadResponse{
		CrawledWebsite: "",
		Undone:         0,
		Processing:     0,
		CrawlerError:   0,
		Saving:         0,
		SendToCrawler:  0,
		Finished:       0,
	}, nil).Times(1)
	crawler.EXPECT().StatusOfTaskQueue(gomock.Any(), gomock.Any()).Return(&crawlerproto.TaskStatusResponse{
		Website:    "",
		Undone:     0,
		Processing: 0,
		Finished:   0,
		Failed:     0,
	}, nil).Times(1)
	sut.Dependency = &SherlockWebServerDependency{
		Crawler: func() crawlerproto.CrawlerService {
			return crawler
		},
		Analyser: func() analyserproto.AnalyserService {
			return analyser
		},
	}

	analyser.EXPECT().StateRPC(gomock.Any(), gomock.Any()).Return(&analyserproto.StateResponse{
		State: analyserproto.AnalyserStateEnum_Clean,
	}, nil)
	crawler.EXPECT().GetState(gomock.Any(), gomock.Any()).Return(nil, &WebserverTestError{})

	sut.GetServiceStatus(c)

	var got gin.H
	err := json.Unmarshal(w.Body.Bytes(), &got)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "200 OK", w.Result().Status)
	assert.Equal(t, "Unknown", fmt.Sprint(got["Crawler"]))
	assert.Equal(t, "Clean", fmt.Sprint(got["Analyser"]))
}

func TestGetServiceStatusAnalyserError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/meta", bytes.NewReader([]byte("{\n        \"operation\": \"Running\",\n        \"target\": \"All\"\n    }")))
	sut := New()
	ctrl := gomock.NewController(t)

	crawler := mockcrawler.NewMockCrawlerService(ctrl)
	analyser := mockanalyser.NewMockAnalyserService(ctrl)
	analyser.EXPECT().WorkloadRPC(gomock.Any(), gomock.Any()).Return(&analyserproto.WorkloadResponse{
		CrawledWebsite: "",
		Undone:         0,
		Processing:     0,
		CrawlerError:   0,
		Saving:         0,
		SendToCrawler:  0,
		Finished:       0,
	}, nil).Times(1)
	crawler.EXPECT().StatusOfTaskQueue(gomock.Any(), gomock.Any()).Return(&crawlerproto.TaskStatusResponse{
		Website:    "",
		Undone:     0,
		Processing: 0,
		Finished:   0,
		Failed:     0,
	}, nil).Times(1)
	sut.Dependency = &SherlockWebServerDependency{
		Crawler: func() crawlerproto.CrawlerService {
			return crawler
		},
		Analyser: func() analyserproto.AnalyserService {
			return analyser
		},
	}

	analyser.EXPECT().StateRPC(gomock.Any(), gomock.Any()).Return(nil, &WebserverTestError{})
	crawler.EXPECT().GetState(gomock.Any(), gomock.Any()).Return(&crawlerproto.StateGetResponse{
		State: crawlerproto.CurrentState_Running,
	}, nil)

	sut.GetServiceStatus(c)

	var got gin.H
	err := json.Unmarshal(w.Body.Bytes(), &got)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "200 OK", w.Result().Status)
	assert.Equal(t, "Running", fmt.Sprint(got["Crawler"]))
	assert.Equal(t, "Unknown", fmt.Sprint(got["Analyser"]))
}

func TestDropGraphTable(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/dropit", nil)
	sut := New()
	ctrl := gomock.NewController(t)

	mockDriver := mocks.NewMockDriver(ctrl)
	mockSession := mocks.NewMockSession(ctrl)

	mockDriver.EXPECT().Session(gomock.Any()).Return(mockSession, nil)

	key := amountofnodes
	value := int64(199)

	mockRecord := mocks.NewMockRecord(ctrl)
	mockResult := mocks.NewMockResult(ctrl)

	mockSession.EXPECT().Run(gomock.Any(), nil).Return(mockResult, nil).Times(1)

	first := mockResult.EXPECT().Next().Return(true).Times(1)
	second := mockResult.EXPECT().Next().Return(false).Times(1)

	gomock.InOrder(first, second)

	mockRecord.EXPECT().Keys().Return([]string{key}).Times(1)
	mockRecord.EXPECT().Get(key).Return(value, true)
	mockResult.EXPECT().Record().Return(mockRecord).Times(2)
	mockSession.EXPECT().Close().Times(1)

	sut.Driver = mockDriver

	sut.DropGraphTable(c)

	var got gin.H
	err := json.Unmarshal(w.Body.Bytes(), &got)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "200 OK", w.Result().Status)
	assert.Equal(t, "Dropped the table.", fmt.Sprint(got["Message"]))
}

func TestDropGraphTableDBError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/dropit", nil)
	sut := New()
	ctrl := gomock.NewController(t)

	mockDriver := mocks.NewMockDriver(ctrl)
	mockSession := mocks.NewMockSession(ctrl)

	mockDriver.EXPECT().Session(gomock.Any()).Return(mockSession, nil)

	key := amountofnodes
	value := int64(199)

	mockRecord := mocks.NewMockRecord(ctrl)
	mockResult := mocks.NewMockResult(ctrl)

	mockSession.EXPECT().Run(gomock.Any(), nil).Return(mockResult, nil).Times(1)

	first := mockResult.EXPECT().Next().Return(true).Times(1)
	second := mockResult.EXPECT().Next().Return(false).Times(1)

	gomock.InOrder(first, second)

	mockRecord.EXPECT().Keys().Return([]string{key}).Times(1)
	mockRecord.EXPECT().Get(key).Return(value, true)
	mockResult.EXPECT().Record().Return(mockRecord).Times(2)
	mockSession.EXPECT().Close().Times(1)

	sut.DropGraphTable(c)

	var got gin.H
	err := json.Unmarshal(w.Body.Bytes(), &got)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "500 Internal Server Error", w.Result().Status)
	assert.Equal(t, "A Problem occurred while trying to drop the Database", fmt.Sprint(got["Message"]))
}

func TestDropGraphTableDriverError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/dropit", nil)
	sut := New()
	ctrl := gomock.NewController(t)

	mockDriver := mocks.NewMockDriver(ctrl)

	mockDriver.EXPECT().Session(gomock.Any()).Return(nil, &WebserverTestError{})
	sut.Driver = mockDriver
	sut.DropGraphTable(c)

	var got gin.H
	err := json.Unmarshal(w.Body.Bytes(), &got)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "500 Internal Server Error", w.Result().Status)
	assert.Equal(t, "A Problem occurred while trying to connect to the Database", fmt.Sprint(got["Message"]))
}
