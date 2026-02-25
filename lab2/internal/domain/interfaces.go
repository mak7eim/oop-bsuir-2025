package domain

type Entity interface {
	GetID() string
}

type FinancialProduct interface {
    GetBalance() float64
    GetCurrency() string
}