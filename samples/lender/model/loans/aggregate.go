package loans

import (
	"fmt"
	"github.com/AltScore/gothic/pkg/ids"

	"github.com/AltScore/gothic/pkg/es"
	"github.com/AltScore/gothic/pkg/es/event"
)

type Event = event.Event

const EntityType = "loans"

type Aggregate struct {
	es.AggregateBase[*LoanView]
}

// New creates a new aggregate with a new Id.
func New() *Aggregate {
	return &Aggregate{
		AggregateBase: es.NewAgg[*LoanView](ids.New(), EntityType, nil, es.WithSnapshot(&LoanView{})),
	}
}

func (a *Aggregate) State() State {
	return a.Snapshot().State
}

func (a *Aggregate) StartFlow(cmd StartFlowCmd) error {
	if a.State() != "" {
		return fmt.Errorf("flow already started")
	}

	return a.Raise(event.For(a, LoanFlowStarted, &FlowStarted{
		ClientID:      cmd.ClientID,
		TransactionID: cmd.TransactionID,
		TotalAmount:   cmd.TotalAmount,
	}))
}

func (a *Aggregate) AcceptTermsAndConditions(
	term int,
	deferredPct Percent,
	acceptConditions bool,
) error {
	if a.State() != Started {
		return fmt.Errorf("flow not started")
	}

	return a.Raise(event.For(a, "termsAndConditionsAccepted", &TermsAndConditionsAccepted{
		Term:             term,
		DeferredPct:      deferredPct,
		AcceptConditions: acceptConditions,
	}))
}

func (a *Aggregate) ConfirmEmail() error {
	if a.Snapshot().State != Accepted {
		return fmt.Errorf("flow not accepted")
	}

	return a.Raise(event.For(a, "emailConfirmed", &EmailConfirmed{}))
}
