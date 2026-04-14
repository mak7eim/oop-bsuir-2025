package domain

const (
	Land  string = "Земля"
	Water string = "Вода"
	Air   string = "Воздух"
)

type Transport struct {
	Name        string
	Type        string
	CostPerKm   float64
	SpeedPerKmH float64
}

func (t *Transport) CalculateDeliveryCost(distance float64) float64 {
	return distance * t.CostPerKm
}

func (t *Transport) CalculateTime(distance float64) float64 {
	return distance / t.SpeedPerKmH
}
