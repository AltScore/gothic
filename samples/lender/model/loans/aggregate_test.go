package loans

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
