package sample1

import (
	"fmt"
	"sync"
	"time"
)

// PriceService is a service that we can use to get prices for the items
// Calls to this service are expensive (they take time)
type PriceService interface {
	GetPriceFor(itemCode string) *PriceModel
}

// TransparentCache is a cache that wraps the actual service
// The cache will remember prices we ask for, so that we don't have to wait on every call
// Cache should only return a price if it is not older than "maxAge", so that we don't get stale prices
type TransparentCache struct {
	actualPriceService PriceService
	maxAge             time.Duration
	Prices             map[string]*PriceModel
}

type PriceModel struct {
	Price *Price
	Error error
}

type Price struct {
	amount    float64
	updatedAt time.Time
}

func NewTransparentCache(actualPriceService PriceService, maxAge time.Duration) *TransparentCache {
	return &TransparentCache{
		actualPriceService: actualPriceService,
		maxAge:             maxAge,
		Prices:             map[string]*PriceModel{},
	}
}

// GetPriceFor gets the price for the item, either from the cache or the actual service if it was not cached or too old
func (c *TransparentCache) GetPriceFor(itemCode string) *PriceModel {
	price, ok := c.Prices[itemCode]
	if ok && c.isLessMaxAge(itemCode) {
		return price
	}

	price = c.actualPriceService.GetPriceFor(itemCode)
	if price.Error != nil {

		return &PriceModel{
			Price: &Price{
				amount: 0, updatedAt: time.Time{},
			},
			Error: fmt.Errorf("getting price from service : %v", price.Error.Error()),
		}
	}
	c.Prices[itemCode] = price
	c.Prices[itemCode].Price.updatedAt = time.Now()
	return price
}

// isLessMaxAge checks that the price was retrieved less than "maxAge" ago
func (c *TransparentCache) isLessMaxAge(itemCode string) bool {
	var isLessMaxAge bool
	now := time.Now()
	updatedDuration := now.Sub(c.Prices[itemCode].Price.updatedAt)

	if updatedDuration < c.maxAge {
		isLessMaxAge = true
	}

	return isLessMaxAge
}

// GetPricesFor gets the prices for several items at once, some might be found in the cache, others might not
// If any of the operations returns an error, it should return an error as well
func (c *TransparentCache) GetPricesFor(itemCodes ...string) ([]*Price, error) {
	pricesBuff := make(chan *PriceModel, len(itemCodes))
	var wg sync.WaitGroup
	for _, itemCode := range itemCodes {
		wg.Add(1)
		go func(id string) {
			pricesBuff <- c.GetPriceFor(id)
			wg.Done()
		}(itemCode)
	}

	wg.Wait()
	close(pricesBuff)

	var prices []*Price
	for pricesDom := range pricesBuff {
		if pricesDom.Error != nil {
			return []*Price{}, pricesDom.Error
		}
		prices = append(prices, pricesDom.Price)
	}

	return prices, nil
}
