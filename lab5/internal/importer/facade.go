package importer

import (
	"fmt"
	"lab5/internal/domain"
	"path/filepath"
	"strings"
)

// DataImporter — фасад для импорта данных (паттерн Фасад)
type DataImporter struct {
	strategies map[string]ImportStrategy
}

func NewDataImporter() *DataImporter {
	return &DataImporter{
		strategies: map[string]ImportStrategy{
			".json": &JSONStrategy{},
			".xml":  &XMLStrategy{},
			".csv":  &CSVStrategy{},
		},
	}
}

// ImportFromFile определяет формат по расширению и применяет нужную стратегию
func (im *DataImporter) ImportFromFile(filePath string) ([]domain.Transport, map[string]CargoInfo, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	strategy, ok := im.strategies[ext]
	if !ok {
		return nil, nil, fmt.Errorf("неподдерживаемый формат файла: %s", ext)
	}

	data, err := readFile(filePath)
	if err != nil {
		return nil, nil, err
	}

	return strategy.Parse(data)
}

// ImportFromBytes — импорт из байт с явным указанием формата
func (im *DataImporter) ImportFromBytes(data []byte, format string) ([]domain.Transport, map[string]CargoInfo, error) {
	var strategy ImportStrategy
	switch strings.ToLower(format) {
	case "json":
		strategy = &JSONStrategy{}
	case "xml":
		strategy = &XMLStrategy{}
	case "csv":
		strategy = &CSVStrategy{}
	default:
		return nil, nil, fmt.Errorf("неподдерживаемый формат: %s", format)
	}

	return strategy.Parse(data)
}
