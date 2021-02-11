package datarail

import "deal-hunter/scrapers"

type DealConsumer interface {
	consume(deal scrapers.Deal)
}

type DataRail struct {
}
