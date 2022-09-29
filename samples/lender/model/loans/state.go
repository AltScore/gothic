package loans

import (
	"github.com/AltScore/gothic/pkg/es/event"
)

type State string

const (
	Started   State = "started"
	Accepted  State = "accepted"
	Confirmed State = "confirmed"
)

type LoanView struct {
	ID               ID
	Version          int
	ClientID         ClientID
	TransactionID    string
	TotalAmount      Money
	State            State
	Term             int
	DeferredPct      Percent
	IsEmailConfirmed bool
}

func (v *LoanView) Apply(event event.IEvent) error {
	switch e := event.Data().(type) {
	case *FlowStarted:
		return e.Apply(v)
	case *TermsAndConditionsAccepted:
		return e.Apply(v)
	case *EmailConfirmed:
		return e.Apply(v)
	}
	return nil
}

func (v *LoanView) SetVersion(version int) {
	v.Version = version
}
