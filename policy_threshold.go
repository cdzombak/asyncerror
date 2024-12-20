package asyncerror

import (
	"log"
	"sync"
	"time"
)

// ThresholdEscalationPolicy will escalate an error if more than ErrorCount errors
// are received within TimeWindow.
type ThresholdEscalationPolicy struct {
	// ErrorCount is the number of errors that must be received within TimeWindow for escalation to occur.
	ErrorCount int
	// TimeWindow is the duration within which ErrorCount errors must be received for escalation to occur.
	TimeWindow time.Duration
	// Name is the name of the policy.
	Name string
	// UniqID is a unique identifier for this policy instance.
	UniqID string
	// LogEvery, if nonzero, will cause the policy to log every Nth received error using log.Println.
	LogEvery int

	lastCompression     time.Time
	skippedSinceLastLog int
	errors              []errorTimeRecord
	mutex               sync.Mutex
}

func (e *ThresholdEscalationPolicy) Close()                    {}
func (e *ThresholdEscalationPolicy) GetDesiredBufferSize() int { return e.ErrorCount * 2 }
func (e *ThresholdEscalationPolicy) GetName() string           { return e.Name }
func (e *ThresholdEscalationPolicy) GetUniqID() string         { return e.UniqID }

func (e *ThresholdEscalationPolicy) Receive(err error) bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if e.LogEvery > 0 {
		e.skippedSinceLastLog++
		if e.skippedSinceLastLog >= e.LogEvery {
			go log.Println(err.Error())
			e.skippedSinceLastLog = 0
		}
	}

	now := time.Now()
	errorsInWindow := 0
	performCompress := now.Sub(e.lastCompression) > e.TimeWindow
	var compressedErrors []errorTimeRecord

	e.errors = append(e.errors, errorTimeRecord{
		At:  now,
		Err: err,
	})

	if performCompress {
		newSliceSize := len(e.errors) / 2
		if newSliceSize <= 1 {
			newSliceSize = 2
		}
		compressedErrors = make([]errorTimeRecord, 0, newSliceSize)
	}

	for _, errRecord := range e.errors {
		if now.Sub(errRecord.At) <= e.TimeWindow {
			errorsInWindow++

			if performCompress {
				compressedErrors = append(compressedErrors, errRecord)
			}
		}
	}

	if performCompress {
		e.lastCompression = now
		e.errors = compressedErrors
	}

	return errorsInWindow >= e.ErrorCount
}

type errorTimeRecord struct {
	At  time.Time
	Err error
}
