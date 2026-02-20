package domain

type Bank struct {
	id string
	name string
	branches []*Branch
}

func NewBank(id, name string) *Bank {
	return &Bank{
		id: id,
		name: name,
		branches: []*Branch{},
	}
}

func (b *Bank) GetID() string {
	return b.id
}

func (b *Bank) GetName() string {
	return b.name
}

func (b *Bank) GetAllBranch() []*Branch {
	result := make([]*Branch, len(b.branches))
	copy(result, b.branches)
	return result
}

func (b *Bank) AddBranch(branch *Branch) {
	b.branches = append(b.branches, branch)
}