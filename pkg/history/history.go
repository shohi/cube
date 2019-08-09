package history

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/shohi/cube/pkg/base"
)

func Write() error {
	h := newHist()
	if h.err != nil {
		return h.err
	}

	h.addNewRecord()

	return h.write()
}

func Read() error {
	lines, err := readAllLines(base.DefaultHistoryPath)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s\n", strings.Join(lines, "\n"))
	return nil
}

func Delete(from int) error {
	if from <= 0 {
		return emptyHistory()
	}

	h := newHist()
	if h.err != nil {
		return h.err
	}

	return h.delete(from)
}

type hist struct {
	records

	lowIndex  int
	highIndex int
	err       error
}

func newHist() hist {
	lines, err := readAllLines(base.DefaultHistoryPath)
	if err != nil {
		return hist{err: err}
	}

	var h hist
	initHist(&h, lines)

	return h
}

func (h hist) write() error {
	content := strings.Join(h.records.lines(), "\n")

	return ioutil.WriteFile(base.DefaultHistoryPath,
		[]byte(content), 0666)
}
func (h *hist) addNewRecord() {
	idx := h.highIndex + 1
	cmd := strings.Join(os.Args, " ")
	line := fmt.Sprintf("%d  %s", idx, cmd)

	h.records = append(h.records, record{index: idx, line: line})
}

func (h *hist) delete(from int) error {
	if from < h.lowIndex {
		return emptyHistory()
	}

	if from > h.highIndex {
		return nil
	}

	// TODO: refactor
	newRs := make(records, 0, len(h.records))

	for _, r := range h.records {
		if r.index < from {
			newRs = append(newRs, r)
		}
	}

	h.records = newRs

	return h.write()
}

func emptyHistory() error {
	file, err := os.Create(base.DefaultHistoryPath)
	if err != nil {
		return err
	}

	if file != nil {
		_ = file.Close()
	}

	return nil
}

func initHist(h *hist, lines []string) {
	if len(lines) == 0 {
		return
	}

	var min int
	var max int
	var initialized bool

	rs := make(records, 0, len(lines))

	for _, line := range lines {
		r, err := newRecord(line)
		if err != nil {
			h.err = err
			return
		}

		rs = append(rs, r)

		if !initialized {
			// skip empty line
			if r.index > 0 {
				min, max = r.index, r.index
				initialized = true
				continue
			}
		}

		if r.index > max {
			max = r.index
		}
		if r.index < min {
			min = r.index
		}
	}

	h.records = rs
	h.lowIndex = min
	h.highIndex = max

	return
}

// https://stackoverflow.com/questions/8757389/reading-file-line-by-line-in-go
func readAllLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := make([]string, 0, 1024)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
