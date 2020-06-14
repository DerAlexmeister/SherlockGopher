package sherlockcrawler

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type logTester struct {

}

func (hf *logTester) Run() {
	stuff, _ := readLines("../info.log")

	start := make([]string, 0)
	end := make([]string, 0)

	for _, v := range stuff {
		parts := strings.Split(v, " ")
		if strings.Contains(v, "STARTGET") {
			start = append(start, parts[2])
		} else {
			end = append(end, parts[2])
		}
	}

	fmt.Println(start)
	fmt.Println(end)

	for _, v := range start {
		if !contains(end, v) {
			fmt.Println(v)
		}
	}

}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
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
