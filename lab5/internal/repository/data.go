package repository

import (
	"encoding/csv"
	"fmt"
	"lab5/internal/domain"
	"os"
	"strconv"
	"strings"
)

type DataLoader struct{}

func NewDataLoader() *DataLoader {
	return &DataLoader{}
}

func (dl *DataLoader) LoadData(filePath string) ([]domain.Transport, map[string]CargoInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка открытия файла: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка чтения CSV: %w", err)
	}

	if len(records) < 2 {
		return nil, nil, fmt.Errorf("файл не содержит данных")
	}

	var transports []domain.Transport
	cargoInfoMap := make(map[string]CargoInfo)

	for i, record := range records[1:] { // Пропускаем заголовок
		if len(record) < 6 {
			continue
		}

		recordType := strings.TrimSpace(record[0])

		switch recordType {
		case "transport":
			transport, err := dl.parseTransport(record)
			if err != nil {
				fmt.Printf("Ошибка парсинга транспорта в строке %d: %v\n", i+2, err)
				continue
			}
			transports = append(transports, transport)

		case "cargo":
			cargoType, info, err := dl.parseCargo(record)
			if err != nil {
				fmt.Printf("Ошибка парсинга груза в строке %d: %v\n", i+2, err)
				continue
			}
			cargoInfoMap[cargoType] = info
		}
	}

	return transports, cargoInfoMap, nil
}

func (dl *DataLoader) parseTransport(record []string) (domain.Transport, error) {
	var transport domain.Transport

	transport.Name = strings.TrimSpace(record[1])

	name := transport.Name
	switch {
	case strings.Contains(name, "Земля"):
		transport.Type = domain.Land
	case strings.Contains(name, "Вода"):
		transport.Type = domain.Water
	case strings.Contains(name, "Воздух"):
		transport.Type = domain.Air
	}

	if len(record) > 4 && record[4] != "" {
		cost, err := strconv.ParseFloat(strings.Replace(record[4], ",", ".", -1), 64)
		if err != nil {
			return transport, fmt.Errorf("ошибка парсинга расхода на км: %w", err)
		}
		transport.CostPerKm = cost
	}

	// Колонка 5: Скорость_км_ч
	if len(record) > 5 && record[5] != "" {
		speed, err := strconv.ParseFloat(strings.Replace(record[5], ",", ".", -1), 64)
		if err != nil {
			return transport, fmt.Errorf("ошибка парсинга скорости: %w", err)
		}
		transport.SpeedPerKmH = speed
	}

	return transport, nil
}

func (dl *DataLoader) parseCargo(record []string) (string, CargoInfo, error) {
	name := strings.TrimSpace(record[1])
	var cargoType string

	switch name {
	case "Электроника":
		cargoType = domain.Electronic
	case "Одежда":
		cargoType = domain.Clothing
	case "Оборудование":
		cargoType = domain.Equipment
	case "Скоропортящиеся продукты":
		cargoType = domain.Perishable
	default:
		return "", CargoInfo{}, fmt.Errorf("неизвестный тип груза: %s", name)
	}

	weight, err := strconv.ParseFloat(strings.Replace(record[2], ",", ".", -1), 64)
	if err != nil {
		return "", CargoInfo{}, fmt.Errorf("ошибка парсинга массы: %w", err)
	}

	transportCost, err := strconv.ParseFloat(strings.Replace(record[3], ",", ".", -1), 64)
	if err != nil {
		return "", CargoInfo{}, fmt.Errorf("ошибка парсинга стоимости перевозки: %w", err)
	}

	return cargoType, CargoInfo{
		WeightPerKg:        weight,
		TransportCostPerKg: transportCost,
	}, nil
}

type CargoInfo struct {
	WeightPerKg        float64
	TransportCostPerKg float64
}
