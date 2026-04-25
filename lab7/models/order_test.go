package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"lab7-oop/models"
)

func TestOrderStatusFlow(t *testing.T) {
	next, ok := models.OrderStatusAccepted.Next()
	assert.True(t, ok)
	assert.Equal(t, models.OrderStatusPreparing, next)

	next, ok = models.OrderStatusOnTheWay.Next()
	assert.True(t, ok)
	assert.Equal(t, models.OrderStatusDelivered, next)

	_, ok = models.OrderStatusDelivered.Next()
	assert.False(t, ok)
}
