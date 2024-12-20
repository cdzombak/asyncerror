package asyncerror

import (
	"fmt"
	"sync"
)

const (
	defaultErrChanBufferSize = 32
)

// NewEscalator returns a new Escalator.
func NewEscalator() Escalator {
	return &escalator{
		errEscalationChan: make(chan error, defaultErrChanBufferSize),
	}
}

// Escalator manages a set of error policies.
// Each policy can consume errors that occur in your program and decide whether to escalate them.
// If any policy decides that an error should be escalated, it will be sent to the EscalationChannel.
// Your application should read from this channel and handle escalated errors (often by logging
// them and/or stopping the program).
type Escalator interface {
	// EscalationChannel returns a channel that will receive errors escalated by this escalator's policies.
	EscalationChannel() chan error

	// RegisterPolicy registers a new policy with this escalator.
	// It returns a channel, to which your application should send errors.
	// The policy must have a UniqID or Name.
	RegisterPolicy(policy Policy) chan error

	// UnregisterPolicy unregisters a policy from this escalator.
	UnregisterPolicy(policy Policy)
}

func (h *escalator) EscalationChannel() chan error {
	return h.errEscalationChan
}

func (h *escalator) RegisterPolicy(policy Policy) chan error {
	h.policiesMutex.Lock()
	defer h.policiesMutex.Unlock()

	if h.policies == nil {
		h.policies = make(map[string]policyRecord)
	}
	policyUid := uidForPolicy(policy)
	if _, ok := h.policies[policyUid]; ok {
		panic(fmt.Sprintf("policy '%s' is already registered", policyUid))
	}

	bufSize := policy.GetDesiredBufferSize()
	if bufSize <= 0 {
		bufSize = defaultErrChanBufferSize
	}
	errorChan := make(chan error, bufSize)
	closeChan := make(chan struct{})

	go func() {
		for {
			select {
			case err := <-errorChan:
				go func() {
					name := policy.GetName()
					uid := policy.GetUniqID()
					if name == "" {
						name = "<unnamed>"
					}
					if policy.Receive(err) {
						if uid != "" {
							h.errEscalationChan <- fmt.Errorf("async error policy '%s' (%s) escalated: %w", name, uid, err)
						} else {
							h.errEscalationChan <- fmt.Errorf("async error policy '%s' escalated: %w", name, err)
						}
					}
				}()
			case <-closeChan:
				return
			}
		}
	}()

	h.policies[policyUid] = policyRecord{
		uid: policyUid,
		closer: func() {
			close(errorChan)
			policy.Close()
			close(closeChan)
		},
	}

	return errorChan
}

func (h *escalator) UnregisterPolicy(policy Policy) {
	h.policiesMutex.Lock()
	defer h.policiesMutex.Unlock()

	policyUid := uidForPolicy(policy)
	if _, ok := h.policies[policyUid]; !ok {
		panic(fmt.Sprintf("policy '%s' is not registered", policyUid))
	}

	h.policies[policyUid].closer()
	delete(h.policies, policyUid)
}

type escalator struct {
	errEscalationChan chan error
	policiesMutex     sync.Mutex
	policies          map[string]policyRecord
}

type policyRecord struct {
	closer func()
	uid    string
}
