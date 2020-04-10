package sherlockanalyser

import (
	"bufio"
	"io/ioutil"
	"os"
	"testing"
)

var testData = []struct {
	in   string
	out  string
	addr string
}{
	{"./testfiles/in1.txt", "./testfiles/out1.txt", "https://github.com/jwalteri/GO-StringUtils"},
	{"./testfiles/in2.txt", "./testfiles/out2.txt", "https://walterj.de/mmix.html"},
}

func TestAnalyser(t *testing.T) {
	for _, tt := range testData {
		t.Run(tt.in, func(t *testing.T) {
			atask := NewTask()

			htmlcode, _ := ioutil.ReadFile(tt.in)
			expected, _ := readLines(tt.out)

			atask.setHTMLCode(string(htmlcode))
			atask.setAddr(tt.addr)
			atask.setTaskID(1)
			atask.Execute()

			if len(expected) != len(atask.getFoundLinks()) {
				t.Errorf("got %d elements, want %d elements", len(atask.foundLinks), len(expected))
			}

			for i, ele := range atask.getFoundLinks() {
				if ele != expected[i] {
					t.Errorf("got %q, want %q", ele, expected[i])
				}
			}
		})
	}
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
