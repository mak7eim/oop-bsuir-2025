package domain

type Employee struct {
	Person
	position string
	employed bool
}

func NewEmployee(id, firstName, lastName, phoneNum, position string) *Employee {
	return &Employee{
		Person: Person{id: id, lastName: lastName, firstName: firstName, phoneNum: phoneNum},
		position: position,
		employed: true,
	}
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