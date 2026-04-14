package importer

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"lab5/internal/domain"
	"strconv"
	"strings"
)

// CargoInfo — временная структура для передачи данных между импортером и фабрикой
type CargoInfo struct {
	WeightPerKg        float64
	TransportCostPerKg float64
}

// ImportStrategy — интерфейс стратегии импорта (паттерн Стратегия)
type ImportStrategy interface {
	Parse(data []byte) ([]domain.Transport, map[string]CargoInfo, error)
}

// --- Общая модель для парсинга (JSON/XML/CSV) ---

type itemRecord struct {
	RecordType      string   `json:"Тип записи" xml:"Тип_записи"`
	Name            string   `json:"Наименование" xml:"Наименование"`
	WeightPerKg     *float64 `json:"Масса_ед_кг" xml:"Масса_ед_кг"`
	TransportCostKg *float64 `json:"Стоимость_перевозки_за_кг" xml:"Стоимость_перевозки_за_кг"`
	CostPerKm       *float64 `json:"Расход_на_км" xml:"Расход_на_км"`
	SpeedKmH        *float64 `json:"Скорость_км_ч" xml:"Скорость_км_ч"`
}

type recordsWrapper struct {
	Items []itemRecord `json:"items" xml:"item"`
}

// --- общая функция маппинга записей ---

func mapRecords(records []itemRecord) ([]domain.Transport, map[string]CargoInfo, error) {
	var transports []domain.Transport
	cargoMap := make(map[string]CargoInfo)

	for _, r := range records {
		switch strings.TrimSpace(r.RecordType) {
		case "transport":
			t, err := mapToTransport(r)
			if err != nil {
				return nil, nil, fmt.Errorf("ошибка парсинга транспорта '%s': %w", r.Name, err)
			}
			transports = append(transports, t)
		case "cargo":
			ct, info, err := mapToCargo(r)
			if err != nil {
				return nil, nil, fmt.Errorf("ошибка парсинга груза '%s': %w", r.Name, err)
			}
			cargoMap[ct] = info
		default:
			continue
		}
	}
	return transports, cargoMap, nil
}

func mapToTransport(r itemRecord) (domain.Transport, error) {
	t := domain.Transport{Name: strings.TrimSpace(r.Name)}

	// Определяем тип по названию
	name := t.Name
	switch {
	case strings.Contains(name, "Земля"):
		t.Type = domain.Land
	case strings.Contains(name, "Вода"):
		t.Type = domain.Water
	case strings.Contains(name, "Воздух"):
		t.Type = domain.Air
	default:
		return t, fmt.Errorf("не удалось определить тип транспорта: %s", name)
	}

	if r.CostPerKm != nil {
		t.CostPerKm = *r.CostPerKm
	}
	if r.SpeedKmH != nil {
		t.SpeedPerKmH = *r.SpeedKmH
	}

	return t, nil
}

func mapToCargo(r itemRecord) (string, CargoInfo, error) {
	var cargoType string
	name := strings.TrimSpace(r.Name)

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

	info := CargoInfo{}
	if r.WeightPerKg != nil {
		info.WeightPerKg = *r.WeightPerKg
	}
	if r.TransportCostKg != nil {
		info.TransportCostPerKg = *r.TransportCostKg
	}

	return cargoType, info, nil
}

// --- JSON стратегия ---

type JSONStrategy struct{}

func (s *JSONStrategy) Parse(data []byte) ([]domain.Transport, map[string]CargoInfo, error) {
	// Пробуем сначала как обертку {"items": [...]}, затем как массив [...]
	var wrapper recordsWrapper
	var records []itemRecord

	if err := json.Unmarshal(data, &wrapper); err == nil && len(wrapper.Items) > 0 {
		records = wrapper.Items
	} else if err := json.Unmarshal(data, &records); err != nil {
		return nil, nil, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	return mapRecords(records)
}

// --- XML стратегия ---

type XMLStrategy struct{}

func (s *XMLStrategy) Parse(data []byte) ([]domain.Transport, map[string]CargoInfo, error) {
	var wrapper recordsWrapper
	if err := xml.Unmarshal(data, &wrapper); err != nil {
		return nil, nil, fmt.Errorf("ошибка парсинга XML: %w", err)
	}
	return mapRecords(wrapper.Items)
}

// --- CSV стратегия ---

type CSVStrategy struct{}

func (s *CSVStrategy) Parse(data []byte) ([]domain.Transport, map[string]CargoInfo, error) {
	reader := csv.NewReader(strings.NewReader(string(data)))
	reader.Comma = ';'
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка чтения CSV: %w", err)
	}

	if len(records) < 2 {
		return nil, nil, fmt.Errorf("файл не содержит данных")
	}

	var items []itemRecord

	for i, row := range records[1:] { // пропускаем заголовок
		if len(row) < 6 {
			continue
		}
		item := parseCSVRow(row)
		if item.RecordType == "" {
			continue
		}
		items = append(items, item)
		_ = i
	}

	return mapRecords(items)
}

func parseCSVRow(row []string) itemRecord {
	item := itemRecord{
		RecordType: strings.TrimSpace(row[0]),
		Name:       strings.TrimSpace(row[1]),
	}

	// Масса_ед_кг
	if v, err := strconv.ParseFloat(strings.Replace(strings.TrimSpace(row[2]), ",", ".", -1), 64); err == nil {
		item.WeightPerKg = &v
	}
	// Стоимость_перевозки_за_кг
	if v, err := strconv.ParseFloat(strings.Replace(strings.TrimSpace(row[3]), ",", ".", -1), 64); err == nil {
		item.TransportCostKg = &v
	}
	// Расход_на_км
	if s := strings.TrimSpace(row[4]); s != "" {
		if v, err := strconv.ParseFloat(strings.Replace(s, ",", ".", -1), 64); err == nil {
			item.CostPerKm = &v
		}
	}
	// Скорость_км_ч
	if s := strings.TrimSpace(row[5]); s != "" {
		if v, err := strconv.ParseFloat(strings.Replace(s, ",", ".", -1), 64); err == nil {
			item.SpeedKmH = &v
		}
	}

	return item
}
