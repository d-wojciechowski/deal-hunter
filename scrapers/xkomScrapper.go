package scrapers

import (
	"encoding/json"
	"github.com/gocolly/colly"
	"github.com/google/logger"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var MoneyRegexp = regexp.MustCompile("[ zł]")

func ScrapXKomGroup(root string) *Deal {
	logger.Info("----------------------------------------------------------------------------------------")
	logger.Infof("Start parsing %s in XKomGroupParser", root)
	logger.Info("New collector init")
	c := colly.NewCollector()
	deal := &Deal{}

	c.OnHTML(".hot-shot", func(e *colly.HTMLElement) {
		pImpression := e.DOM.Find(".product-impression")
		deal.Name = pImpression.Find(".product-name").Text()
		logger.Infof("Parsed name :%s", deal.Name)
		deal.Link = root + "goracy_strzal"
		logger.Infof("Parsed link :%s", deal.Link)
		deal.ImgLink, _ = pImpression.Find("img").Attr("src")
		deal.ImgLink = strings.Replace(deal.ImgLink, "?filters=grayscale", "", -1)
		logger.Infof("Parsed img link :%s", deal.ImgLink)
		priceDiv := e.DOM.Find(".price")

		oldPrice := MoneyRegexp.ReplaceAllString(priceDiv.Find(".old-price").Text(), "")
		oldPrice = strings.Replace(oldPrice, ",", ".", 1)
		deal.OldPrice, _ = strconv.ParseFloat(oldPrice, 64)
		logger.Infof("Parsed old price :%0.2f", deal.OldPrice)

		newPrice := MoneyRegexp.ReplaceAllString(priceDiv.Find(".new-price").Text(), "")
		newPrice = strings.Replace(newPrice, ",", ".", 1)
		deal.NewPrice, _ = strconv.ParseFloat(newPrice, 64)
		logger.Infof("Parsed new price :%0.2f", deal.NewPrice)

		countDiv := e.DOM.Find(".count")
		deal.Left, _ = strconv.ParseInt(countDiv.Find(".pull-Left > .gs-quantity").Text(), 10, 64)
		logger.Infof("Parsed left count :%d", deal.Left)
		deal.Sold, _ = strconv.ParseInt(countDiv.Find(".pull-right > .gs-quantity").Text(), 10, 64)
		logger.Infof("Parsed sold count :%d", deal.Sold)

		deal.Start = getStartDate()
		logger.Infof("Parsed start date :%s", deal.Start)
		deal.End = getEndDate()
		logger.Infof("Parsed end date :%s", deal.End)
	})

	c.OnRequest(func(r *colly.Request) {
		logger.Infof("Visiting %s", r.URL.String())
	})

	err := c.Visit(root)
	if err != nil {
		logger.Errorf("Could not parse %s", root)
		logger.Error(err.Error())
	}
	marshall, _ := json.MarshalIndent(deal, "", "\t")
	logger.Infof("Scrapped object:\n%s", string(marshall))
	logger.Info("----------------------------------------------------------------------------------------")
	return deal
}

func getStartDate() time.Time {
	now := time.Now()
	if now.Hour() > 9 && now.Hour() < 22 {
		return time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, now.Location())
	}
	return time.Date(now.Year(), now.Month(), now.Day(), 22, 0, 0, 0, now.Location())
}

func getEndDate() time.Time {
	now := time.Now()
	if now.Hour() > 9 && now.Hour() < 22 {
		return time.Date(now.Year(), now.Month(), now.Day(), 21, 59, 59, 0, now.Location())
	}
	return time.Date(now.Year(), now.Month(), now.Day()+1, 9, 59, 59, 0, now.Location())
}
