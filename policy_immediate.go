package asyncerror

import "log"

// ImmediateEscalationPolicy will escalate every error it receives.
type ImmediateEscalationPolicy struct {
	// Name is the name of the policy.
	Name string

	// UniqID is a unique identifier for this policy instance.
	UniqID string

	// Whether to log the error using log.Println.
	Log bool
}

func (i *ImmediateEscalationPolicy) Close()                    {}
func (i *ImmediateEscalationPolicy) GetDesiredBufferSize() int { return 1 }
func (i *ImmediateEscalationPolicy) GetName() string           { return i.Name }
func (i *ImmediateEscalationPolicy) GetUniqID() string         { return i.UniqID }

func (i *ImmediateEscalationPolicy) Receive(err error) bool {
	if i.Log {
		log.Println(err.Error())
	}
	return true
}
