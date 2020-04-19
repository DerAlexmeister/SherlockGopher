package sherlockanalyser

import (
	"bufio"
	"fmt"
	test "github.com/ob-algdatii-20ss/SherlockGopher/analyser/sherlockanalyser/test"
	crawlerproto "github.com/ob-algdatii-20ss/SherlockGopher/sherlockcrawler/proto/crawlertoanalyser"
	"io/ioutil"
	"os"
	"testing"
)

var testData = []struct {
	in   string
	out  string
	addr string
}{
	{"./test/in1.txt", "./test/out1.txt", "https://github.com/jwalteri/GO-StringUtils"},
	{"./test/in2.txt", "./test/out2.txt", "https://walterj.de/mmix.html"},
}

func TestAnalyser(t *testing.T) {
	for _, tt := range testData {
		t.Run(tt.in, func(t *testing.T) {
			htmlcode, _ := ioutil.ReadFile(tt.in)
			expected, _ := readLines(tt.out)

			cdata := CrawlerData{
				taskId:            1,
				addr:              tt.addr,
				taskError:         nil,
				responseHeader:    nil,
				responseBodyBytes: htmlcode,
				statusCode:        200,
				responseTime:      0,
			}

			task := injectDependencies(NewTask(cdata))

			task.Execute(nil)

			if len(expected) != len(task.FoundLinks()) {
				t.Errorf("got %d elements, want %d elements", len(task.foundLinks), len(expected))
			}

			for i, ele := range task.FoundLinks() {
				if ele != expected[i] {
					t.Errorf("got %q, want %q", ele, expected[i])
				}
			}

			fmt.Print("Parser:")
			fmt.Println(task.parserTime)

			fmt.Print("Anaylser:")
			fmt.Println(task.analyserTime)
		})
	}
}

func TestErrorCase(t *testing.T) {
	cData := CrawlerData{
		taskId:            1,
		addr:              "www.error.err",
		taskError:         fmt.Errorf("test"),
		responseHeader:    nil,
		responseBodyBytes: []byte(""),
		statusCode:        200,
		responseTime:      0,
	}

	task := NewTask(cData)
	task = injectDependencies(task)

	task.Execute(nil)

	if len(task.FoundLinks()) != 0 {
		t.Errorf("task was worngly analyzed")
	}

}

func TestGetterSetter(t *testing.T) {
	state := PROCESSING
	html := "html"
	workAdrr := "wAdr"
	cData := &CrawlerData{}
	foundLinks := make([]string, 0)
	rootAdr := "rAdr"
	var aTime int64 = 1
	var id uint64 = 1
	linkTags := make(map[string]string)
	var pTime int64 = 1

	task := analyserTaskRequest{}
	task.SetState(state)
	task.SetHtml(html)
	task.SetWorkAddr(workAdrr)
	task.SetCrawlerData(cData)
	task.SetFoundLinks(foundLinks)
	task.SetRootAddr(rootAdr)
	task.SetAnalyserTime(aTime)
	task.SetId(id)
	task.SetLinkTags(linkTags)
	task.SetParserTime(pTime)

	if len(task.LinkTags()) != len(linkTags) {
		t.Errorf("linkTags is wrong")
	}
	if task.ParserTime() != pTime {
		t.Errorf("pTime is wrong")
	}
	if task.Id() != id {
		t.Errorf("id is wrong")
	}
	if task.AnalyserTime() != aTime {
		t.Errorf("aTime is wrong")
	}
	if task.RootAddr() != rootAdr {
		t.Errorf("rootAdr is wrong")
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
	if task.Html() != html {
		t.Errorf("html is wrong")
	}
	if task.State() != state {
		t.Errorf("state is wrong")
	}
}

func injectDependencies(task analyserTaskRequest) analyserTaskRequest {
	neo4j := test.GetNeo4jSessionInstance()
	crawler := test.GetAnalyserInterfaceServiceInstance()
	crawlerFunc := func() crawlerproto.AnalyserInterfaceService {
		return crawler
	}
	task.InjectDependency(&AnalyserDependency{
		Neo4J: &neo4j,
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
