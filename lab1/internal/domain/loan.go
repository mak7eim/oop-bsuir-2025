package domain


type Loan struct {
	id         string
	clientID   string
	money      float64
	currency   string
}

func NewLoan(id, clientID string, currency string, money float64, termMonths int) *Loan {
	return &Loan{
		id: id,
		clientID: clientID,
		currency: currency,
		money: money,
	}
}

func (l *Loan) GetID() string {
	return l.id
}

func (l *Loan) GetClientID() string {
	return l.clientID
}

func (l *Loan) GetRemainingDebt() float64 {
	return l.money
}

func (l *Loan) MakePayment(payment float64) {
	l.money -= payment
}
