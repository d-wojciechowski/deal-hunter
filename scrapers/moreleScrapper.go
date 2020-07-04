package scrapers

import (
	"github.com/gocolly/colly"
	"github.com/google/logger"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type MoreleScrapper struct {
	URL *url.URL
	collyDriven
}

var MoneyRegexp = regexp.MustCompile("[  zł]")
var NumberRegexp = regexp.MustCompile("[0-9]+")

func (scrapper *MoreleScrapper) Scrap() *Deal {
	c := scrapper.collyDriven.newColly(scrapper.URL.String())
	deal := &Deal{SiteName: "morele"}

	c.OnHTML(".home-sections-promotion", func(e *colly.HTMLElement) {
		href := e.DOM.Find(".prom-box-top").Find("a")
		deal.Name, _ = href.Attr("title")
		deal.Link, _ = href.Attr("href")
		deal.ImgLink, _ = href.Find("img").Attr("src")

		priceDiv := e.DOM.Find(".promo-box-price")

		oldPrice := MoneyRegexp.ReplaceAllString(priceDiv.Find(".promo-box-old-price").Text(), "")
		oldPrice = strings.Replace(oldPrice, ",", ".", -1)
		deal.OldPrice, _ = strconv.ParseFloat(oldPrice, 64)

		newPrice := MoneyRegexp.ReplaceAllString(priceDiv.Find(".promo-box-new-price").Text(), "")
		newPrice = strings.Replace(newPrice, ",", ".", -1)
		deal.NewPrice, _ = strconv.ParseFloat(newPrice, 64)

		countDiv := e.DOM.Find(".status-box-labels")

		leftTest := countDiv.Find(".status-box-was").Text()
		deal.Left, _ = strconv.ParseInt(NumberRegexp.FindString(leftTest), 10, 64)

		soldTest := countDiv.Find(".status-box-expired").Text()
		deal.Sold, _ = strconv.ParseInt(NumberRegexp.FindString(soldTest), 10, 64)

		deal.Code = e.DOM.Find(".promo-box-code-value").Text()

		deal.Start = time.Now()
		deal.End = getMoreleEndDate(e.DOM.Find(".promo-box-countdown").Attr("data-date-to"))
	})

	err := c.Visit(scrapper.URL.String())
	if err != nil {
		logger.Errorf("Could not parse %s", scrapper.URL.String())
		logger.Error(err.Error())
	}

	scrapper.collyDriven.logDeal(deal)
	return deal
}

func getMoreleEndDate(date string, exists bool) time.Time {
	value, _ := time.Parse("2006-01-02 15:04:05", date)
	return value
}
