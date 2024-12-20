package asyncerror

import "fmt"

// Policy defines how to handle errors that occur in an asynchronous context.
type Policy interface {
	// Close is called when the policy is unregistered.
	// Use it for e.g. cleanup of state & resources.
	Close()

	// GetDesiredBufferSize returns the desired size for this policy's error channel buffer.
	GetDesiredBufferSize() int

	// GetName returns a human-readable name for the policy.
	// This may be (but is not required to be) the same for all instances of this Policy.
	GetName() string

	// GetUniqID returns a unique identifier for the policy.
	// This must be unique across all policies in use.
	GetUniqID() string

	// Receive is called by the Escalator when an error occurs. Your policy should return
	// true if the error should be escalated further.
	Receive(err error) bool
}

func uidForPolicy(policy Policy) string {
	if policy.GetUniqID() != "" {
		return policy.GetUniqID()
	}
	if policy.GetName() != "" {
		return policy.GetName()
	}
	panic(fmt.Sprintf("policy has no Name or UniqID: %v", policy))
}
