package history

import (
	"strconv"
	"strings"
)

type record struct {
	index int
	line  string
}

type records []record

func (rs records) lines() []string {
	if len(rs) == 0 {
		return nil
	}

	lines := make([]string, 0, len(rs))
	for _, r := range rs {
		lines = append(lines, r.line)
	}

	return lines
}

func newRecord(line string) (record, error) {
	r := record{line: line}
	if len(line) == 0 {
		return r, nil
	}

	tokens := strings.Fields(line)
	idx, err := strconv.ParseInt(tokens[0], 10, 32)
	if err != nil {
		return r, err
	}

	r.index = int(idx)

	return r, nil
}
