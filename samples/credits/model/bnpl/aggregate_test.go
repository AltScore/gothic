package bnpl

import "testing"

func Test_Reify_BNPL(t *testing.T) {

	bnpl1, _ := Reify(nil)

	err := bnpl1.StartFlow(StartFlowCmd{
		ClientID:      "client1",
		TransactionID: "transaction1",
		TotalAmount:   100,
	})
	if err != nil {
		return
	}

	err = bnpl1.AcceptTermsAndConditions(30, 75, true)

}
