/*
Package es implements a common support for Event Sourcing.

To use, you should extend your Entity Aggregate from the AggregateBase struct.

Example:

	package loans

	import (
		"fmt"

		"github.com/AltScore/gothic/pkg/es"
	)

	const EntityType = "loans"

	type Aggregate struct {
		es.AggregateBase[ID, Snapshot]
	}

Your events should extend the Event struct:

	type flowStarted struct {
		es.Metadata
		ClientID      ClientID
		TransactionID string
		TotalAmount   Money
	}


Package samples contains sample code for using gothic.
*/
package es
