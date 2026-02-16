package domain

type Transaction struct {
	id          string
	fromAccount string
	toAccount   string
	money       float64
	currency    string
	status      string //в ожидании, выполнено, ошибка
}

func NewTransaction(id, fromAccount, toAccount, currency string, money float64) *Transaction {
	return &Transaction{
		id:          id,
		fromAccount: fromAccount,
		toAccount:   toAccount,
		currency:    currency,
		money:       money,
		status:      "в ожидании",
	}
}

func (t *Transaction) GetID() string {
	return t.id
}

func (t *Transaction) GetMoney() float64 {
	return t.money
}

func (t *Transaction) GetStatus() string {
	return t.status
}

func (t *Transaction) Complete() {
	t.status = "выполнено"
}

func (t *Transaction) Fail() {
	t.status = "ошибка"
}