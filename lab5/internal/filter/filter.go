package filter

import "lab5/internal/domain"

// Filter — интерфейс фильтра
type Filter interface {
	Apply(results []domain.DeliveryResult) []domain.DeliveryResult
}

// TransportNameFilter — фильтр по названию транспорта (содержит подстроку)
type TransportNameFilter struct {
	Substring string
}

func (f *TransportNameFilter) Apply(results []domain.DeliveryResult) []domain.DeliveryResult {
	if f.Substring == "" {
		return results
	}
	var out []domain.DeliveryResult
	for _, r := range results {
		if containsIgnoreCase(r.Transport.Name, f.Substring) {
			out = append(out, r)
		}
	}
	return out
}

// MaxPriceFilter — фильтр по максимальной итоговой стоимости
type MaxPriceFilter struct {
	MaxPrice float64
}

func (f *MaxPriceFilter) Apply(results []domain.DeliveryResult) []domain.DeliveryResult {
	var out []domain.DeliveryResult
	for _, r := range results {
		if r.TotalCost <= f.MaxPrice {
			out = append(out, r)
		}
	}
	return out
}

func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && searchIgnoreCase(s, substr)
}

func searchIgnoreCase(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if equalIgnoreCase(s[i:i+len(substr)], substr) {
			return true
		}
	}
	return false
}

func equalIgnoreCase(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if toLower(a[i]) != toLower(b[i]) {
			return false
		}
	}
	return true
}

func toLower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + 32
	}
	return c
}
