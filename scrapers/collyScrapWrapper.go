package scrapers

import (
	"github.com/gocolly/colly"
	"github.com/google/logger"
)

type collyWrapper struct {
	Link     string
	Selector string
	Consumer func(deal *Deal, element *colly.HTMLElement)
	hostname string
	deal     *Deal
	coly     *colly.Collector
}

func (wrapper *collyWrapper) init() {
	wrapper.coly = colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.97 Safari/537.36 OPR/65.0.3467.48"),
	)
	wrapper.coly.OnRequest(func(r *colly.Request) {
		logger.Infof("Visiting %s", r.URL.String())
	})
	wrapper.coly.OnHTML(wrapper.Selector, func(element *colly.HTMLElement) {
		wrapper.Consumer(wrapper.deal, element)
	})
	wrapper.deal = &Deal{SiteName: wrapper.hostname}
}

func (wrapper *collyWrapper) execute() *Deal {
	err := wrapper.coly.Visit(wrapper.Link)
	if err != nil {
		logger.Errorf("Could not parse %s", wrapper.Link)
		logger.Error(err.Error())
	}
	return wrapper.deal
}
