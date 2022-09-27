package bnpl

import (
	"fmt"

	"github.com/AltScore/gothic/pkg/es"
)

const EntityType = "bnpl"

type Aggregate struct {
	base es.AggregateBase[ID, Snapshot]
}

// Reify recreates an aggregate from a list of events stored to its current state.
func Reify(previousEvents []Event) (*Aggregate, error) {
	a := Aggregate{
		base: es.NewAgg[ID, Snapshot](NewId(), EntityType, previousEvents),
	}

	return &a, a.base.Replay()
}

func (a *Aggregate) State() State {
	return a.base.Snapshot().State
}

func (a *Aggregate) StartFlow(cmd StartFlowCmd) error {
	if a.State() != "" {
		return fmt.Errorf("flow already started")
	}

	return a.base.Raise(
		FlowStarted{
			Metadata:      a.base.NewMetadata("FlowStarted"),
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

	return a.base.Raise(
		TermsAndConditionsAccepted{
			Metadata:         a.base.NewMetadata("termsAndConditionsAccepted"),
			Term:             term,
			DeferredPct:      deferredPct,
			AcceptConditions: acceptConditions,
		},
	)
}

func (a *Aggregate) ConfirmEmail() error {
	if a.base.Snapshot().State != Accepted {
		return fmt.Errorf("flow not accepted")
	}

	return a.base.Raise(
		EmailConfirmed{
			Metadata: a.base.NewMetadata("emailConfirmed"),
		},
	)
}
