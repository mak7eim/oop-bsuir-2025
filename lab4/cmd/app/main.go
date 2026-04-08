package main

import (
	"fmt"
	"lab4/internal/domain"
	"lab4/internal/factory"
	"lab4/internal/repository"
	"lab4/internal/service"
)

func main() {
	loader := repository.NewDataLoader()
	transports, cargoInfo, err := loader.LoadData("../../materials/logistic.csv")
	if err != nil {
		panic(err)
	}

	factoryCargoInfo := make(map[string]factory.CargoInfo)
	for k, v := range cargoInfo {
		factoryCargoInfo[k] = factory.CargoInfo{
			WeightPerKg: v.WeightPerKg,
			CostPerKg:   v.CostPerKg,
		}
	}

	provider := factory.NewLogisticsFactory(transports, factoryCargoInfo)

	calculator := service.NewDeliveryCalculator(provider)

	requests := []service.CargoRequest{
		{Type: domain.Electronic, Quantity: 5},
		{Type: domain.Clothing, Quantity: 10},
	}

	result, err := calculator.CalculateDelivery(requests, domain.Land, 150.0)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}
