package service

import (
	"lab5/internal/domain"
	"lab5/internal/factory"
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
	transportName string, // новый параметр
	distance float64,
) (*domain.DeliveryResult, error) {

	f, err := dc.provider.GetFactory(transportType)
	if err != nil {
		return nil, err
	}

	transport, err := f.CreateTransport(transportName) // передаем имя
	if err != nil {
		return nil, err
	}

	var cargos []domain.Cargo
	var totalTransportCost float64

	for _, req := range cargoItems {
		cargo, err := f.CreateCargo(req.Type, req.Quantity)
		if err != nil {
			return nil, err
		}
		cargos = append(cargos, *cargo)
		totalTransportCost += cargo.TotalTransportCost()
	}

	deliveryCost := transport.CalculateDeliveryCost(distance)
	totalCost := totalTransportCost + deliveryCost

	return &domain.DeliveryResult{
		CargoList:     cargos,
		Transport:     transport,
		Distance:      distance,
		TransportCost: totalTransportCost,
		DeliveryCost:  deliveryCost,
		TotalCost:     totalCost,
		DeliveryHours: transport.CalculateTime(distance),
	}, nil
}

// CalculateAllTransports — если транспорт не указан, возвращает результаты для всех видов транспорта
func (dc *DeliveryCalculator) CalculateAllTransports(
	cargoItems []CargoRequest,
	distance float64,
) ([]domain.DeliveryResult, error) {

	transportTypes := []string{domain.Land, domain.Water, domain.Air}
	var results []domain.DeliveryResult

	for _, tt := range transportTypes {
		f, err := dc.provider.GetFactory(tt)
		if err != nil {
			continue
		}

		// Получаем все транспорты этого типа (через фабрику с пустым именем)
		transport, err := f.CreateTransport("")
		if err != nil {
			continue
		}

		var cargos []domain.Cargo
		var totalTransportCost float64

		for _, req := range cargoItems {
			cargo, err := f.CreateCargo(req.Type, req.Quantity)
			if err != nil {
				return nil, err
			}
			cargos = append(cargos, *cargo)
			totalTransportCost += cargo.TotalTransportCost()
		}

		deliveryCost := transport.CalculateDeliveryCost(distance)
		totalCost := totalTransportCost + deliveryCost

		results = append(results, domain.DeliveryResult{
			CargoList:     cargos,
			Transport:     transport,
			Distance:      distance,
			TransportCost: totalTransportCost,
			DeliveryCost:  deliveryCost,
			TotalCost:     totalCost,
			DeliveryHours: transport.CalculateTime(distance),
		})
	}

	return results, nil
}
