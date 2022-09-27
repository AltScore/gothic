package loans

import (
	"github.com/AltScore/gothic/pkg/es"
)

type ClientID string
type Money float64
type Percent float64

type Event = es.Event[ID, Snapshot]

type StartFlowCmd struct {
	ClientID      ClientID
	TransactionID string
	TotalAmount   Money
}

type FlowStarted struct {
	es.Metadata[ID]
	ClientID      ClientID
	TransactionID string
	TotalAmount   Money
}

func (f FlowStarted) Apply(state *Snapshot) error {
	state.ClientID = f.ClientID
	state.TransactionID = f.TransactionID
	state.TotalAmount = f.TotalAmount
	state.State = Started
	return nil
}

type TermsAndConditionsAccepted struct {
	es.Metadata[ID]

	Term             int
	DeferredPct      Percent
	AcceptConditions bool
}

func (t TermsAndConditionsAccepted) Apply(snapshot *Snapshot) error {
	snapshot.State = Accepted
	snapshot.Term = t.Term
	snapshot.DeferredPct = t.DeferredPct
	return nil
}

type EmailConfirmed struct {
	es.Metadata[ID]
}

func (e EmailConfirmed) Apply(snapshot *Snapshot) error {
	snapshot.IsEmailConfirmed = true
	snapshot.State = Confirmed
	return nil
}

type PhoneConfirmed struct {
	es.Metadata[ID]
}
