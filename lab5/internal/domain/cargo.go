package domain

const (
	Electronic string = "Электроника"
	Clothing   string = "Одежда"
	Equipment  string = "Оборудование"
	Perishable string = "Скоропортящиеся продукты"
)

type Cargo struct {
	Name               string
	Type               string
	WeightPerUnit      float64
	TransportCostPerKg float64
	Quantity           int
}

func (c *Cargo) TotalWeight() float64 {
	return c.WeightPerUnit * float64(c.Quantity)
}

func (c *Cargo) TotalTransportCost() float64 {
	return c.TotalWeight() * c.TransportCostPerKg
}
