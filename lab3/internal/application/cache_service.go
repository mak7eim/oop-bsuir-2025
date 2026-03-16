package application

import (
	"lab3/internal/domain"
)

type LocalCacheService struct {
	data map[string]*domain.Order
}

func NewLocalCacheService() *LocalCacheService {
	return &LocalCacheService{
		data: make(map[string]*domain.Order),
	}
}

func (lcs *LocalCacheService) TryAddOrder(order *domain.Order) bool {
	if _, exists := lcs.data[order.ID]; exists {
		return false
	}

	lcs.data[order.ID] = cloneOrder(order)

	return true
}

func (lcs *LocalCacheService) RemoveOrder(id string) {
	delete(lcs.data, id)
}

func (lcs *LocalCacheService) FindOrder(id string) (*domain.Order, bool) {
	order, exist := lcs.data[id]
	if !exist {
		return nil, false
	}

	return cloneOrder(order), true
}

func cloneOrder(order *domain.Order) *domain.Order {
	if order == nil {
		return nil
	}

	orderCopy := *order
	if len(order.Items) > 0 {
		orderCopy.Items = make([]domain.Item, len(order.Items))
		copy(orderCopy.Items, order.Items)
	}

	return &orderCopy
}
