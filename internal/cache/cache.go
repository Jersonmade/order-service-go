package cache

import (
	"sync"

	"github.com/Jersonmade/test-wb-project/internal/model"
)

type OrderCache struct {
	data sync.Map
}

func NewOrderCache() *OrderCache {
	return &OrderCache{}
}

func (c *OrderCache) Get(orderUID string) (model.Order, bool) {
	value, ok := c.data.Load(orderUID)
	if !ok {
		return model.Order{}, false
	}

	order, ok := value.(model.Order)
	return order, ok
}

func (c *OrderCache) Set(orderUID string, order model.Order)  {
	c.data.Store(orderUID, order)
}