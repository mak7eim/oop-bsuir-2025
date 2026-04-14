package sorter

import (
	"lab5/internal/domain"
	"sort"
)

// SortStrategy — интерфейс стратегии сортировки (паттерн Стратегия)
type SortStrategy interface {
	Sort(results []domain.DeliveryResult)
}

// ByTransportName — сортировка по названию транспорта (по возрастанию)
type ByTransportName struct{}

func (s *ByTransportName) Sort(results []domain.DeliveryResult) {
	sort.SliceStable(results, func(i, j int) bool {
		return results[i].Transport.Name < results[j].Transport.Name
	})
}

// ByPrice — сортировка по итоговой стоимости (по возрастанию)
type ByPrice struct{}

func (s *ByPrice) Sort(results []domain.DeliveryResult) {
	sort.SliceStable(results, func(i, j int) bool {
		return results[i].TotalCost < results[j].TotalCost
	})
}

// BySpeed — сортировка по скорости транспорта (по убыванию — быстрее = выше)
type BySpeed struct{}

func (s *BySpeed) Sort(results []domain.DeliveryResult) {
	sort.SliceStable(results, func(i, j int) bool {
		return results[i].Transport.SpeedPerKmH > results[j].Transport.SpeedPerKmH
	})
}

// CompositeSorter — применяет несколько стратегий последовательно
type CompositeSorter struct {
	strategies []SortStrategy
}

func NewCompositeSorter(strategies ...SortStrategy) *CompositeSorter {
	return &CompositeSorter{strategies: strategies}
}

func (c *CompositeSorter) Sort(results []domain.DeliveryResult) {
	// Применяем в обратном порядке для стабильной сортировки
	// (последняя стратегия будет иметь наивысший приоритет)
	for i := len(c.strategies) - 1; i >= 0; i-- {
		c.strategies[i].Sort(results)
	}
}
