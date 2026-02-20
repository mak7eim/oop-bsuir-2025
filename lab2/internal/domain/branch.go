package domain

type Branch struct {
	id      string
	name    string
	address string
	phone   string
	opened  bool
	atms []*ATM
	employees []*Employee
}

func NewBranch(id, address, phone, name string) *Branch {
	return &Branch{
		id:      id,
		name:    name,
		address: address,
		phone:   phone,
		opened:  true,
		atms: []*ATM{},
		employees: []*Employee{},
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

func (b *Branch) GetAllATM() []*ATM {
	result := make([]*ATM, len(b.atms))
	copy(result, b.atms)
	return result
}

func (b *Branch) GetAllEmployee() []*Employee {
	result := make([]*Employee, len(b.employees))
	copy(result, b.employees)
	return result
}

func (b *Branch) AddATM(atm *ATM) {
	b.atms = append(b.atms, atm)
}

func (b *Branch) AddEmployee(employee *Employee) {
	b.employees = append(b.employees, employee)
}


