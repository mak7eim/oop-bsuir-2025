package domain

// бфнковский счет
type Account struct {
	clientID string
	number   string
	currency string
	balance  float64
	cards    []*Card
}

func NewAccount(number string, clientID string, currency string, balance float64) *Account {
	return &Account{
		number:   number,
		clientID: clientID,
		currency: currency,
		balance:  balance,
		cards: []*Card{},
	}
}

func (a *Account) GetNumber() string {
	return a.number
}

func (a *Account) GetClientID() string {
	return a.clientID
}

func (a *Account) GetBalance() float64 {
	return a.balance
}

func (a *Account) GetCurrency() string {
	return a.currency
}

func (a *Account) Deposit(money float64) {
	a.balance += money
}

func (a *Account) WithdrawMoney(money float64) bool {
	if a.balance >= money {
		a.balance -= money
		return true
	}

	return false
}

func (a *Account) GetAllCard() []*Card {
	result := make([]*Card, len(a.cards))
	copy(result, a.cards)
	return result
}

func (a *Account) AddCard(card *Card) {
	a.cards = append(a.cards, card)
}
