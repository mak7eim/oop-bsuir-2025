package factory

import "lab5/internal/domain"

type TransportFactory interface {
	CreateTransport(transportName string) (domain.Transport, error)
	CreateCargo(cargoType string, quantity int) (*domain.Cargo, error)
}

type CargoInfo struct {
	WeightPerKg        float64
	TransportCostPerKg float64
}
