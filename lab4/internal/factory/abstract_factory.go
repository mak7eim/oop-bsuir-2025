package factory

import "lab4/internal/domain"

type TransportFactory interface {
	CreateTransport() (domain.Transport, error)
	CreateCargo(cargoType string, quantity int) (*domain.Cargo, error)
}
