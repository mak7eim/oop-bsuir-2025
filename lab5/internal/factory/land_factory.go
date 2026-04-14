package factory

import (
	"fmt"
	"lab5/internal/domain"
)

type LandFactory struct {
	transports    []domain.Transport
	cargoBaseInfo map[string]CargoInfo
}

func NewLandFactory(transports []domain.Transport, cargoInfo map[string]CargoInfo) *LandFactory {
	return &LandFactory{
		transports:    transports,
		cargoBaseInfo: cargoInfo,
	}
}

func (f *LandFactory) CreateTransport(transportName string) (domain.Transport, error) {
	if transportName == "" {
		for _, t := range f.transports {
			if t.Type == domain.Land {
				return t, nil
			}
		}
	}

	for _, t := range f.transports {
		if t.Type == domain.Land && t.Name == transportName {
			return t, nil
		}
	}

	return domain.Transport{}, fmt.Errorf("наземный транспорт '%s' не найден", transportName)
}

func (f *LandFactory) CreateCargo(cargoType string, quantity int) (*domain.Cargo, error) {
	info, ok := f.cargoBaseInfo[cargoType]
	if !ok {
		return nil, fmt.Errorf("неизвестный тип груза: %s", cargoType)
	}

	var name string
	switch cargoType {
	case domain.Electronic:
		name = "Электроника"
	case domain.Clothing:
		name = "Одежда"
	case domain.Equipment:
		name = "Оборудование"
	case domain.Perishable:
		name = "Скоропортящиеся продукты"
	default:
		name = cargoType
	}

	return &domain.Cargo{
		Name:               name,
		Type:               cargoType,
		WeightPerUnit:      info.WeightPerKg,
		TransportCostPerKg: info.TransportCostPerKg,
		Quantity:           quantity,
	}, nil
}
