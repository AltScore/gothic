package loans

type State string

const (
	Started   State = "started"
	Accepted  State = "accepted"
	Confirmed State = "confirmed"
)

type LoanView struct {
	ID               ID
	Version          int
	ClientID         ClientID
	TransactionID    string
	TotalAmount      Money
	State            State
	Term             int
	DeferredPct      Percent
	IsEmailConfirmed bool
}

func (s LoanView) SetVersion(version int) {
	s.Version = version
}
