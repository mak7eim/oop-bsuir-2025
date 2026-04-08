package service

import (
	"lab4/internal/domain"
	"lab4/internal/factory"
)

type CargoRequest struct {
	Type     string
	Quantity int
}

type DeliveryCalculator struct {
	provider *factory.LogisticsFactory
}

func NewDeliveryCalculator(provider *factory.LogisticsFactory) *DeliveryCalculator {
	return &DeliveryCalculator{provider: provider}
}

func (dc *DeliveryCalculator) CalculateDelivery(
	cargoItems []CargoRequest,
	transportType string,
	distance float64,
) (*domain.DeliveryResult, error) {

	f, err := dc.provider.GetFactory(transportType)
	if err != nil {
		return nil, err
	}

	transport, err := f.CreateTransport()
	if err != nil {
		return nil, err
	}

	var cargos []domain.Cargo
	var totalCargoCost float64

	for _, req := range cargoItems {
		cargo, err := f.CreateCargo(req.Type, req.Quantity)
		if err != nil {
			return nil, err
		}
		cargos = append(cargos, *cargo)
		totalCargoCost += cargo.TotalCost()
	}

	deliveryCost := transport.CalculateDeliveryCost(distance)

	return &domain.DeliveryResult{
		CargoList:     cargos,
		Transport:     transport,
		Distance:      distance,
		CargoCost:     totalCargoCost,
		DeliveryCost:  deliveryCost,
		TotalCost:     totalCargoCost + deliveryCost,
		DeliveryHours: transport.CalculateTime(distance),
	}, nil
}
