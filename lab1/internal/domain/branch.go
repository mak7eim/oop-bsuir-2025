package domain

type Branch struct {
	id      string
	name    string
	address string
	phone   string
	opened  bool
}

func NewBranch(id, address, phone, name string) *Branch {
	return &Branch{
		id:      id,
		name:    name,
		address: address,
		phone:   phone,
		opened:  true,
	}
}

func (b *Branch) GetID() string {
	return b.id
}

func (b *Branch) GetAddress() string {
	return b.address
}

func (b *Branch) GetPhone() string {
	return b.phone
}

func (b *Branch) IsOpen() bool {
	return b.opened
}

func (b *Branch) UpdatePhone(newPhone string) {
	b.phone = newPhone
}

func (b *Branch) SetOpen() {
	b.opened = true
}

func (b *Branch) SetClose() {
	b.opened = false
}
