package loans

const (
	LoanFlowStarted    = "loan.flow.started"
	EmailConfirmedType = "loan.email_confirmed"
)

type ClientID string
type Money float64
type Percent float64

type StartFlowCmd struct {
	ClientID      ClientID
	TransactionID string
	TotalAmount   Money
}

type FlowStarted struct {
	ClientID      ClientID
	TransactionID string
	TotalAmount   Money
}

func (f FlowStarted) Apply(state *LoanView) error {
	state.ClientID = f.ClientID
	state.TransactionID = f.TransactionID
	state.TotalAmount = f.TotalAmount
	state.State = Started
	return nil
}

type TermsAndConditionsAccepted struct {
	Term             int
	DeferredPct      Percent
	AcceptConditions bool
}

func (t TermsAndConditionsAccepted) Apply(snapshot *LoanView) error {
	snapshot.State = Accepted
	snapshot.Term = t.Term
	snapshot.DeferredPct = t.DeferredPct
	return nil
}

type EmailConfirmed struct {
}

func (e EmailConfirmed) Apply(snapshot *LoanView) error {
	snapshot.IsEmailConfirmed = true
	snapshot.State = Confirmed
	return nil
}
