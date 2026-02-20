package domain

type Deposit struct {
	id       string
	clientID string
	money    float64
	rate     float64
	currency string
}

func NewDeposit(id, clientID string, money float64, currency string, rate float64) *Deposit {
	return &Deposit{
		id:       id,
		clientID: clientID,
		money:    money,
		currency: currency,
		rate:     rate,
	}
}

func (d *Deposit) GetID() string {
	return d.id
}

func (d *Deposit) GetBalance() float64 {
	return d.money
}

func (d *Deposit) CalculateInterest() float64 {
	return d.money * d.rate / 100
}

