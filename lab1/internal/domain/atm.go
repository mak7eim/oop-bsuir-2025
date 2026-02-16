package domain

type ATM struct {
	id      string
	balance float64
	branchID string 
}

func NewATM(id string, balance float64, branchID string) *ATM {
	return &ATM{
		id:       id,
		balance:  balance,
		branchID: branchID,
	}
}

func (a *ATM) GetID() string {
	return a.id
}

func (a *ATM) GetBalance() float64 {
	return a.balance
}

func (a *ATM) GetBranchID() string {
	return a.branchID
}

func (a *ATM) WithdrawMoney(money float64) bool {
	if a.balance >= money {
		a.balance -= money
		return true
	}
	return false
}