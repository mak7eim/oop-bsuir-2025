package domain

type Person struct {
	id        string
	firstName string
	lastName  string
	phoneNum  string
}

func (p *Person) GetID() string {
	return p.id
}

func (p *Person) GetFullName() string {
	return p.lastName + " " + p.firstName
}

func (p *Person) GetPhoneNum() string {
	return p.phoneNum
}

func (c *Client) SetPhone(newPhoneNum string) {
	c.phoneNum = newPhoneNum
}
