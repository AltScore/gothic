package loans

import (
	"fmt"

	"github.com/AltScore/gothic/pkg/es"
)

const EntityType = "loans"

type Aggregate struct {
	es.AggregateBase[ID, LoanView]
}

// New creates a new aggregate with a new ID.
func New() *Aggregate {
	return &Aggregate{
		AggregateBase: es.NewAgg[ID, LoanView](NewId(), EntityType, nil, updateVersion),
	}
}

// Reify recreates an aggregate from a list of events stored to its current state.
func Reify(previousEvents []Event) (*Aggregate, error) {
	if len(previousEvents) == 0 {
		return nil, fmt.Errorf("no events to rebuild from")
	}

	a := Aggregate{
		AggregateBase: es.NewAgg[ID, LoanView](previousEvents[0].EntityID(), EntityType, previousEvents, updateVersion),
	}

	return &a, a.Replay()
}

func updateVersion(view *LoanView, version int) {
	view.Version = version
}

func (a *Aggregate) State() State {
	return a.Snapshot().State
}

func (a *Aggregate) StartFlow(cmd StartFlowCmd) error {
	if a.State() != "" {
		return fmt.Errorf("flow already started")
	}

	return a.Raise(
		FlowStarted{
			Metadata:      a.NewMetadata("FlowStarted"),
			ClientID:      cmd.ClientID,
			TransactionID: cmd.TransactionID,
			TotalAmount:   cmd.TotalAmount,
		},
	)
}

func (a *Aggregate) AcceptTermsAndConditions(
	term int,
	deferredPct Percent,
	acceptConditions bool,
) error {
	if a.State() != Started {
		return fmt.Errorf("flow not started")
	}

	return a.Raise(
		TermsAndConditionsAccepted{
			Metadata:         a.NewMetadata("termsAndConditionsAccepted"),
			Term:             term,
			DeferredPct:      deferredPct,
			AcceptConditions: acceptConditions,
		},
	)
}

func (a *Aggregate) ConfirmEmail() error {
	if a.Snapshot().State != Accepted {
		return fmt.Errorf("flow not accepted")
	}

	return a.Raise(
		EmailConfirmed{
			Metadata: a.NewMetadata("emailConfirmed"),
		},
	)
}
