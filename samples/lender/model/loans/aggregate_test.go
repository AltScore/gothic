package loans

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Reify_Loan(t *testing.T) {
	// GIVEN a new Loan Aggregate
	loan := New()

	err := loan.StartFlow(StartFlowCmd{
		ClientID:      "client1",
		TransactionID: "transaction1",
		TotalAmount:   100,
	})

	assert.NoError(t, err)

	err = loan.AcceptTermsAndConditions(30, 75, true)

	assert.NoError(t, err)
}
