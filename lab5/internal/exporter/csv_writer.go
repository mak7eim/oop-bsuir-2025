package exporter

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"lab5/internal/domain"
)

type CSVWriter struct{}

func (w *CSVWriter) Write(results []domain.DeliveryResult) ([]byte, string, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Comma = ';'

	header := []string{
		"Транспорт", "Тип транспорта", "Расстояние_км",
		"Стоимость_перевозки_грузов", "Накладные_расходы",
		"Итоговая_стоимость", "Время_в_пути_ч",
	}
	if err := writer.Write(header); err != nil {
		return nil, "", err
	}

	for _, r := range results {
		row := []string{
			r.Transport.Name,
			r.Transport.Type,
			fmt.Sprintf("%.1f", r.Distance),
			fmt.Sprintf("%.2f", r.TransportCost),
			fmt.Sprintf("%.2f", r.DeliveryCost),
			fmt.Sprintf("%.2f", r.TotalCost),
			fmt.Sprintf("%.2f", r.DeliveryHours),
		}
		if err := writer.Write(row); err != nil {
			return nil, "", err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, "", err
	}

	return buf.Bytes(), ".csv", nil
}
