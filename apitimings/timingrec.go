package apitimings

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var ErrNotTimingRecord = errors.New("Not an API timing record")

type APITimingRec struct {
	MsgId          int
	SourceName     string
	SourceHost     string
	SourceCategory string
	Message        string
}

func NewAPITimingRec(raw string) (*APITimingRec, error) {
	parts := strings.Split(raw, ",")

	id, err := strconv.Atoi(strings.Replace(parts[0], "\"", "", -1))
	if err != nil {
		return nil, err
	}

	return &APITimingRec{
		MsgId:          id,
		SourceName:     parts[1],
		SourceHost:     parts[2],
		SourceCategory: parts[3],
		Message:        extractMessage(raw),
	}, nil
}

func extractMessage(raw string) string {
	parts := strings.Split(raw, ",")
	m := strings.Join(parts[11:], ",")
	m = m[strings.Index(m, "{"):]
	m = m[:len(m)-1]
	m = strings.Replace(m, "\n", "", -1)
	m = strings.Replace(m, "\"\"", "\"", -1)
	return m
}

func unquote(s string) string {
	return s[1 : len(s)-1]
}

func formatPGTime(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}

func (at *APITimingRec) CallRecord() (string, error) {
	var callRecord endToEndTimer

	err := json.Unmarshal([]byte(at.Message), &callRecord)
	if err != nil {
		return "", err
	}

	if callRecord.Name == "" || callRecord.TxnId == "" {
		return "", ErrNotTimingRecord
	}

	sub := callRecord.Tags["sub"]
	aud := callRecord.Tags["aud"]

	return fmt.Sprintf("%s|%s|%t|%s|%s|%s|%s|%s|%d",
		formatPGTime(callRecord.LoggingTimestamp),
		callRecord.TxnId,
		callRecord.Error != "",
		unquote(at.SourceHost),
		unquote(at.SourceCategory),
		callRecord.Name,
		sub,
		aud,
		callRecord.Duration.Nanoseconds()/100000), nil
}

func (at *APITimingRec) ServiceCalls() ([]string, error) {
	var callRecord endToEndTimer
	var serviceCalls []string

	err := json.Unmarshal([]byte(at.Message), &callRecord)
	if err != nil {
		return serviceCalls, err
	}

	for _, c := range callRecord.Contributors {
		if len(c.ServiceCalls) > 0 {
			for _, sc := range c.ServiceCalls {
				sctxt := fmt.Sprintf("%s|%s|%t|%s|%s|%d",
					formatPGTime(callRecord.LoggingTimestamp),
					callRecord.TxnId,
					sc.Error != "",
					sc.Name,
					sc.Endpoint,
					sc.Duration/100000,
				)

				serviceCalls = append(serviceCalls, sctxt)
			}
		}
	}

	return serviceCalls, nil
}
