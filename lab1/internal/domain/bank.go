package domain

type Bank struct {
	id string
	name string
}

func NewBank(id, name string) *Bank {
	return &Bank{
		id: id,
		name: name,
	}
}

func (b *Bank) GetID() string {
	return b.id
}

func (b *Bank) GetName() string {
	return b.name
}