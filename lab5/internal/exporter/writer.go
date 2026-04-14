package exporter

import "lab5/internal/domain"

// ExportWriter — базовый интерфейс для записи результата (паттерн Декоратор)
type ExportWriter interface {
	Write(results []domain.DeliveryResult) ([]byte, string, error) // данные, расширение, ошибка
}
