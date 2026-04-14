package exporter

import (
	"encoding/json"
	"lab5/internal/domain"
)

type JSONWriter struct{}

func (w *JSONWriter) Write(results []domain.DeliveryResult) ([]byte, string, error) {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return nil, "", err
	}
	return data, ".json", nil
}
