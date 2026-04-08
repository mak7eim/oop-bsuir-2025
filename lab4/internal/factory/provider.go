package factory

import (
	"fmt"
	"lab4/internal/domain"
)

type LogisticsFactory struct {
	landFactory  TransportFactory
	waterFactory TransportFactory
	airFactory   TransportFactory
}

func NewLogisticsFactory(transports []domain.Transport, cargoInfo map[string]CargoInfo) *LogisticsFactory {
	return &LogisticsFactory{
		landFactory:  NewLandFactory(transports, cargoInfo),
		waterFactory: NewWaterFactory(transports, cargoInfo),
		airFactory:   NewAirFactory(transports, cargoInfo),
	}
}

func (p *LogisticsFactory) GetFactory(transportType string) (TransportFactory, error) {
	switch transportType {
	case domain.Land:
		return p.landFactory, nil
	case domain.Water:
		return p.waterFactory, nil
	case domain.Air:
		return p.airFactory, nil
	default:
		return nil, fmt.Errorf("неизвестный тип транспорта: %s", transportType)
	}
}
