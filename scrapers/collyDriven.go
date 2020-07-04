package scrapers

import (
	"encoding/json"
	"github.com/gocolly/colly"
	"github.com/google/logger"
	"reflect"
)

type collyDriven struct {
}

func (object *collyDriven) newColly(link string) (c *colly.Collector) {
	logger.Info("----------------------------------------------------------------------------------------")
	logger.Infof("Start parsing %s in %s", link, reflect.TypeOf(object).Name())
	logger.Info("New collector init")
	c = colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.97 Safari/537.36 OPR/65.0.3467.48"),
	)
	c.OnRequest(func(r *colly.Request) {
		logger.Infof("Visiting %s", r.URL.String())
	})
	return
}

func (object *collyDriven) logDeal(deal *Deal) {
	marshall, _ := json.MarshalIndent(deal, "", "\t")
	logger.Infof("Scrapped object:\n%s", string(marshall))
	logger.Info("----------------------------------------------------------------------------------------")
}
