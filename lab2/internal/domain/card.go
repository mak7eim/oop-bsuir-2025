package domain

import "time"

type Card struct {
	cardNum     string
	accountNum  string //номер счета
	expiryMonth int
	expiryYear  int
	cvv         int
}

func NewCard(cardNum, accountNum string, expiryMonth, expityYear, cvv int) *Card {
	return &Card{
		cardNum:     cardNum,
		accountNum:  accountNum,
		expiryMonth: expiryMonth,
		expiryYear:  expityYear,
		cvv:         cvv,
	}
}

func (c *Card) GetCardNum() string {
	return c.cardNum
}

func (c *Card) GetAccountNum() string {
	return c.accountNum
}

func (c *Card) IsExpired() bool {
	return c.expiryYear < time.Now().Year() || (c.expiryYear == time.Now().Year() && c.expiryMonth <= int(time.Now().Month()))
}
