package scrapers

import (
	"fmt"
	"github.com/gocolly/colly"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var MoneyRegexp = regexp.MustCompile("[ zł]")

func ScrapXKomGroup(root string) *Deal {
	// Instantiate default collector
	c := colly.NewCollector()
	deal := &Deal{}

	c.OnHTML(".hot-shot", func(e *colly.HTMLElement) {
		pImpression := e.DOM.Find(".product-impression")
		deal.Name = pImpression.Find(".product-name").Text()
		deal.Link = root + "goracy_strzal"
		deal.ImgLink, _ = pImpression.Find("img").Attr("src")
		deal.ImgLink = strings.Replace(deal.ImgLink, "?filters=grayscale", "", -1)
		priceDiv := e.DOM.Find(".price")

		oldPrice := MoneyRegexp.ReplaceAllString(priceDiv.Find(".old-price").Text(), "")
		oldPrice = strings.Replace(oldPrice, ",", ".", 1)
		deal.OldPrice, _ = strconv.ParseFloat(oldPrice, 64)

		newPrice := MoneyRegexp.ReplaceAllString(priceDiv.Find(".new-price").Text(), "")
		newPrice = strings.Replace(newPrice, ",", ".", 1)
		deal.NewPrice, _ = strconv.ParseFloat(newPrice, 64)

		countDiv := e.DOM.Find(".count")
		deal.Left, _ = strconv.ParseInt(countDiv.Find(".pull-Left > .gs-quantity").Text(), 10, 64)
		deal.Sold, _ = strconv.ParseInt(countDiv.Find(".pull-right > .gs-quantity").Text(), 10, 64)

		deal.Start = getStartDate()
		deal.End = getEndDate()
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit(root)
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
