package apitimings

import (
	"sync"
	"time"
)

type serviceCall struct {
	Name          string
	Endpoint      string
	Duration      time.Duration
	Error         string
	errorReported bool
	start         time.Time
}

type contributor struct {
	Name          string
	Duration      time.Duration
	Error         string
	errorReported bool
	start         time.Time
	ServiceCalls  []*serviceCall
}

type endToEndTimer struct {
	sync.RWMutex
	Name             string
	Tags             map[string]string
	Duration         time.Duration
	LoggingTimestamp time.Time `json:"time"`
	TxnId            string
	Contributors     []*contributor
	ErrorFree        bool
	Error            string
	errorReported    bool
	start            time.Time
}
