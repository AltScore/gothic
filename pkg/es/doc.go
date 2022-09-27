/*
Package es implements a common support for Event Sourcing.

To use, you should extend your Entity Aggregate from the AggregateBase struct.

Example:

	package bnpl

	import (
		"fmt"

		"github.com/AltScore/gothic.git/pkg/es"
	)

	const EntityType = "bnpl"

	type Aggregate struct {
		es.AggregateBase[ID, Snapshot]
	}

Your events should extend the Event struct:

	type flowStarted struct {
		es.Metadata[ID]
		ClientID      ClientID
		TransactionID string
		TotalAmount   Money
	}


Package samples contains sample code for using gothic.
*/
package es
