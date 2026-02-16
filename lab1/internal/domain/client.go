package domain

type Client struct {
	id        string
	firstName string
	lastName  string
	phoneNum  string
	isActive  bool
}

func NewClient(id, firstName, lastName, phone string) *Client {
	return &Client{
		id:        id,
		firstName: firstName,
		lastName:  lastName,
		phoneNum:  phone,
		isActive:  true,
	}
}

func (c *Client) GetID() string {
	return c.id
}

func (c *Client) GetFullName() string {
	return c.firstName + " " + c.lastName
}

func (c *Client) IsActive() bool {
	return c.isActive
}

func (c *Client) DisActive() {
	c.isActive = false
}

func (c *Client) SetPhone(newPhoneNum string) {
	c.phoneNum = newPhoneNum
}
