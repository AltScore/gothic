package loans

import (
	"testing"

	"github.com/AltScore/gothic/pkg/es/event"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewLoanHasVersion0(t *testing.T) {
	// GIVEN a new Loan Aggregate
	loan := New()

	assert.Equal(t, 0, loan.Version())
}

func Test_NewLoanHasNoEvents(t *testing.T) {
	// GIVEN a new Loan Aggregate
	loan := New()

	assert.Empty(t, loan.GetNewEvents())
}

func Test_StartFlow_has_an_event(t *testing.T) {
	// GIVEN a new Loan Aggregate
	loan := New()

	err := loan.StartFlow(StartFlowCmd{
		ClientID:      "client1",
		TransactionID: "transaction1",
		TotalAmount:   100,
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, len(loan.GetNewEvents()))
	assert.Equal(t, 1, loan.Version())
	assert.Equal(t, 1, loan.Snapshot().Version)
}

func Test_StartFlow_has_an_several_events(t *testing.T) {
	// GIVEN a new Loan Aggregate
	loan := New()

	err := loan.StartFlow(StartFlowCmd{
		ClientID:      "client1",
		TransactionID: "transaction1",
		TotalAmount:   100,
	})

	assert.NoError(t, err)

	err = loan.AcceptTermsAndConditions(30, 75, true)
	require.NoError(t, err)

	assert.Equal(t, 2, len(loan.GetNewEvents()))
	assert.Equal(t, 2, loan.Version())
	assert.Equal(t, 2, loan.Snapshot().Version)
}

func TestCan_use_entry(t *testing.T) {
	loan := New()

	err := loan.StartFlow(StartFlowCmd{
		ClientID:      "client1",
		TransactionID: "transaction1",
		TotalAmount:   100,
	})

	assert.NoError(t, err)

	events := loan.GetNewEvents()

	entry := event.From(events[0])

	assert.Equal(t, LoanFlowStarted, entry.Name)
}
