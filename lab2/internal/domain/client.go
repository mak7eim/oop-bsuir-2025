package domain

type Client struct {
	Person
	isActive bool
}

func NewClient(id, firstName, lastName, phoneNum string) *Client {
	return &Client{
		Person:   Person{id: id, firstName: firstName, lastName: lastName, phoneNum: phoneNum},
		isActive: true,
	}
}

func (c *Client) IsActive() bool {
	return c.isActive
}

func (c *Client) DisActive() {
	c.isActive = false
}
