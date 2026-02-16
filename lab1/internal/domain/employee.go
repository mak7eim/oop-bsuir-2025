package domain

type Employee struct {
	id       string
	fullName string
	position string
	employed bool
}

func NewEmployee(id, fullName, position string) *Employee {
	return &Employee{
		id:       id,
		fullName: fullName,
		position: position,
		employed: true,
	}
}

func (e *Employee) GetID() string {
	return e.id
}

func (e *Employee) GetFullName() string {
	return e.fullName
}

func (e *Employee) GetPosition() string {
	return e.position
}

func (e Employee) IsEmployed() bool {
	return e.employed
}

func (e *Employee) StartWork() {
	e.employed = true
}

func (e *Employee) StopWork() {
	e.employed = false
}
