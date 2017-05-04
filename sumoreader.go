package sumoreader

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type SumoReader struct {
	scanner      *bufio.Scanner
	text         string
	linesScanned int
}

func NewSumoReader(reader io.Reader) (*SumoReader, error) {

	scanner := bufio.NewScanner(reader)
	//Skip first line (headers)
	scanner.Scan()

	return &SumoReader{
		scanner:      scanner,
		linesScanned: 1,
	}, nil
}

func (sr *SumoReader) Scan() bool {

	var buffer bytes.Buffer
	defer func() { sr.text = buffer.String() }()

	sr.linesScanned++
	if !sr.scanner.Scan() {
		return false
	}

	line := sr.scanner.Text()
	buffer.WriteString(line)
	if isFullRecord(line) {
		return true
	}

	for scanDone := false; !scanDone; {
		sr.linesScanned++
		more := sr.scanner.Scan()
		if !more {
			return false
		}

		line := sr.scanner.Text()
		buffer.WriteString(line)
		if recordDone(line) {
			scanDone = true
		}
	}

	return true
}

func (sr *SumoReader) Text() string {
	return sr.text
}

func (sr *SumoReader) Err() error {
	return sr.scanner.Err()
}

func isFullRecord(line string) bool {
	parts := strings.Split(line, ".")

	if len(parts) < 13 {
		return false
	}

	return strings.LastIndex(parts[len(parts)-1], "\"") == len(parts[len(parts)-1])-1
}

func recordDone(line string) bool {
	return strings.LastIndex(line, "\"") == len(line)-1
}

func (sr *SumoReader) Lines() int {
	return sr.linesScanned
}
