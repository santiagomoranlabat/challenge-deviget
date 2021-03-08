package sample1

import (
	"fmt"
	"time"
)

// PriceService is a service that we can use to get prices for the items
// Calls to this service are expensive (they take time)
type PriceService interface {
	GetPriceFor(itemCode string) (*Price, error)
}

// TransparentCache is a cache that wraps the actual service
// The cache will remember prices we ask for, so that we don't have to wait on every call
// Cache should only return a price if it is not older than "maxAge", so that we don't get stale prices
type TransparentCache struct {
	actualPriceService PriceService
	maxAge             time.Duration
	Prices             map[string]*Price
}

type Price struct {
	amount    float64
	updatedAt time.Time
}

func NewTransparentCache(actualPriceService PriceService, maxAge time.Duration) *TransparentCache {
	return &TransparentCache{
		actualPriceService: actualPriceService,
		maxAge:             maxAge,
		Prices:             map[string]*Price{},
	}
}

// GetPriceFor gets the price for the item, either from the cache or the actual service if it was not cached or too old
func (c *TransparentCache) GetPriceFor(itemCode string) (*Price, error) {
	price, ok := c.Prices[itemCode]
	if ok && c.isLessMaxAge(itemCode) {
		return price, nil
	}

	price, err := c.actualPriceService.GetPriceFor(itemCode)
	if err != nil {
		return &Price{
			amount:    0,
			updatedAt: time.Now(),
		}, fmt.Errorf("getting amount from service : %v", err.Error())
	}
	c.Prices[itemCode] = price
	c.Prices[itemCode].updatedAt = time.Now()
	return price, nil
}

// isLessMaxAge checks that the price was retrieved less than "maxAge" ago
func (c *TransparentCache) isLessMaxAge(itemCode string) bool {
	var isLessMaxAge bool
	now := time.Now()
	updatedDuration := now.Sub(c.Prices[itemCode].updatedAt)

	if updatedDuration < c.maxAge {
		isLessMaxAge = true
	}

	return isLessMaxAge
}

// GetPricesFor gets the prices for several items at once, some might be found in the cache, others might not
// If any of the operations returns an error, it should return an error as well
func (c *TransparentCache) GetPricesFor(itemCodes ...string) ([]float64, error) {
	results := []float64{}
	for _, itemCode := range itemCodes {
		// TODO: parallelize this, it can be optimized to not make the calls to the external service sequentially
		price, err := c.GetPriceFor(itemCode)
		if err != nil {
			return []float64{}, err
		}
		results = append(results, price.amount)
	}
	return results, nil
}
