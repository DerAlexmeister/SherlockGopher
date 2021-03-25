package sherlockanalyser

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/DerAlexx/SherlockGopher/analyser/sherlockanalyser/test"
	crawlerproto "github.com/DerAlexx/SherlockGopher/sherlockcrawler/proto"
	neo "github.com/DerAlexx/SherlockGopher/sherlockneo"
	"github.com/golang/mock/gomock"
	jp "github.com/jpillora/go-tld"
)

var testData = []struct {
	in    string
	out   string
	addr  string
	count int
}{
	{"./test/in1.txt", "./test/out1.txt", "https://github.com/jwalteri/GO-StringUtils", 43},
	{"./test/in2.txt", "./test/out2.txt", "https://walterj.de/mmix.html", 1},
	{"./test/in3.txt", "./test/out3.txt", "https://www.bbc.co.uk", 2},
	{"./test/in4.txt", "./test/out4.txt", "https://www.bbc.co.uk/news/uk-52493500", 63},
}

func TestAnalyser(t *testing.T) {
	for _, tt := range testData {
		t.Run(tt.in, func(t *testing.T) {
			htmlCode, _ := ioutil.ReadFile(tt.in)
			expected, _ := readLines(tt.out)

			header := http.Header{}
			cdata := CrawlerData{
				taskID:            1,
				addr:              tt.addr,
				taskError:         errors.New(""),
				responseHeader:    &header,
				responseBodyBytes: htmlCode,
				statusCode:        200,
				responseTime:      0,
			}
			mockCtrl := gomock.NewController(t)
			mockNeoSaver := NewMockneoSaverInterface(mockCtrl)
			mockNeoSaver.EXPECT().Save(gomock.Any()).Times(tt.count)
			//mockNeoSaver.EXPECT().Contains(gomock.Any()).MinTimes(tt.count)
			mockNeoSaver.EXPECT().GetSession().Return(nil)
			mockNeoSaver.EXPECT().Contains(gomock.Any()).Return(expected).MinTimes(tt.count)
			task := injectDependencies(NewTask(&cdata))
			cache := NewAnalyserCache()
			task.cache = &cache
			(*task).SetSaver(mockNeoSaver)
			task.Execute(nil)

			if len(expected) != len(task.FoundLinks()) {
				t.Errorf("got %d elements, want %d elements", len(task.foundLinks), len(expected))
			}

			for i, ele := range task.FoundLinks() {
				ele = strings.ReplaceAll(ele, "%2F", "/")
				if ele != expected[i] {
					t.Errorf("got %q, want %q", ele, expected[i])
				}
			}

			if task.FileType() != neo.HTML {
				t.Errorf("filetype is wrong! Got %q, want %q", task.FileType(), neo.HTML)
			}
		})
	}

	fmt.Println("ALL GOOD!")
}

var testFileData = []struct {
	addr     string
	fileType neo.FileType
}{
	{"www.walterj.de/main.html", neo.HTML},
	{"www.walterj.de/main.js", neo.Javascript},
	{"www.walterj.de/main.png", neo.Image},
	{"www.walterj.de/main.jpg", neo.Image},
	{"www.walterj.de/main.jepg", neo.Image},
	{"www.walterj.de/main.gif", neo.Image},
	{"www.walterj.de/main.css", neo.CSS},
	{"www.walterj.de/main.unknown", neo.HTML},
}

func TestExtractFileType(t *testing.T) {

	for _, tt := range testFileData {
		t.Run(tt.addr, func(t *testing.T) {
			domainInfo, _ := jp.Parse(tt.addr)

			task := AnalyserTaskRequest{}
			if task.extractFileType(domainInfo) != tt.fileType {
				t.Fatal("Filetyp is wrong")
			}
		})
	}
}

func TestErrorCase(t *testing.T) {
	cData := CrawlerData{
		taskID:            1,
		addr:              "www.error.err",
		taskError:         fmt.Errorf("test"),
		responseHeader:    nil,
		responseBodyBytes: []byte(""),
		statusCode:        200,
		responseTime:      0,
	}

	mockCtrl := gomock.NewController(t)
	mockNeoSaver := NewMockneoSaverInterface(mockCtrl)
	mockNeoSaver.EXPECT().Save(gomock.Any()).MinTimes(1)
	mockNeoSaver.EXPECT().GetSession().Return(nil)
	mockNeoSaver.EXPECT().Contains(gomock.Any()).Return(make([]string, 0)).MinTimes(1)

	task := injectDependencies(NewTask(&cData))
	(*task).SetSaver(mockNeoSaver)
	task.Execute(nil)

	if task.State() != FINISHED {
		t.Fatal("Error task failed")
	}
}

func TestGetterSetter(t *testing.T) {
	state := PROCESSING
	html := neo.HTML
	workAdrr := "wAdr"
	fileTyp := neo.HTML
	cData := &CrawlerData{}
	foundLinks := make([]string, 0)
	var aTime int64 = 1
	var id uint64 = 1
	var pTime int64 = 1

	task := AnalyserTaskRequest{}
	task.SetState(state)
	task.SetHTML(string(html))
	task.SetWorkAddr(workAdrr)
	task.SetCrawlerData(cData)
	task.SetFoundLinks(foundLinks)
	task.SetAnalyserTime(aTime)
	task.SetID(id)
	task.SetParserTime(pTime)
	task.SetFileType()

	if task.FileType() != fileTyp {
		t.Errorf("fileTyp is wrong")
	}
	if task.ParserTime() != pTime {
		t.Errorf("pTime is wrong")
	}
	if task.getID() != id {
		t.Errorf("id is wrong")
	}
	if task.AnalyserTime() != aTime {
		t.Errorf("aTime is wrong")
	}
	if len(task.FoundLinks()) != len(foundLinks) {
		t.Errorf("foundLinks is wrong")
	}
	if task.CrawlerData() != cData {
		t.Errorf("cData is wrong")
	}
	if task.WorkAddr() != workAdrr {
		t.Errorf("workAdrr is wrong")
	}
	if task.HTML() != string(neo.HTML) {
		t.Errorf("html is wrong")
	}
	if task.State() != state {
		t.Errorf("state is wrong")
	}
}

func injectDependencies(task *AnalyserTaskRequest) *AnalyserTaskRequest {
	neo4j, _ := neo.GetNewDatabaseConnection()
	crawler := test.GetAnalyserInterfaceServiceInstance()
	crawlerFunc := func() crawlerproto.CrawlerService {
		return crawler
	}
	task.InjectDependency(&AnalyserDependency{
		Neo4J:   &neo4j,
		Crawler: crawlerFunc,
	})

	return task
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}
