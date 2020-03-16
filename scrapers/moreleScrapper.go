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

var NumberRegexp = regexp.MustCompile("[0-9]+")

func ScrapMorele() *Deal {
	logger.Info("----------------------------------------------------------------------------------------")
	logger.Infof("Start parsing in MoreleScrapper")
	logger.Info("New collector init")
	c := colly.NewCollector()
	deal := &Deal{}

	c.OnHTML(".home-sections-promotion", func(e *colly.HTMLElement) {
		deal.SiteName = "morele"
		href := e.DOM.Find(".prom-box-top").Find("a")
		deal.Name, _ = href.Attr("title")
		logger.Infof("Parsed name :%s", deal.Name)
		deal.Link, _ = href.Attr("href")
		logger.Infof("Parsed link :%s", deal.Link)
		deal.ImgLink, _ = href.Find("img").Attr("src")
		logger.Infof("Parsed img link :%s", deal.ImgLink)
		priceDiv := e.DOM.Find(".promo-box-price")

		oldPrice := MoneyRegexp.ReplaceAllString(priceDiv.Find(".promo-box-old-price").Text(), "")
		oldPrice = strings.Replace(oldPrice, ",", ".", -1)
		deal.OldPrice, _ = strconv.ParseFloat(oldPrice, 64)
		logger.Infof("Parsed old price :%0.2f", deal.OldPrice)

		newPrice := MoneyRegexp.ReplaceAllString(priceDiv.Find(".promo-box-new-price").Text(), "")
		newPrice = strings.Replace(newPrice, ",", ".", -1)
		deal.NewPrice, _ = strconv.ParseFloat(newPrice, 64)
		logger.Infof("Parsed new price :%0.2f", deal.NewPrice)

		countDiv := e.DOM.Find(".status-box-labels")
		leftTest := countDiv.Find(".status-box-was").Text()
		deal.Left, _ = strconv.ParseInt(NumberRegexp.FindString(leftTest), 10, 64)
		logger.Infof("Parsed left count :%d", deal.Left)

		soldTest := countDiv.Find(".status-box-expired").Text()
		deal.Sold, _ = strconv.ParseInt(NumberRegexp.FindString(soldTest), 10, 64)
		logger.Infof("Parsed sold count :%d", deal.Sold)

		deal.Code = e.DOM.Find(".promo-box-code-value").Text()
		logger.Infof("Parsed discount code :%s", deal.Code)

		deal.Start = time.Now()
		logger.Infof("Parsed start date :%s", deal.Start)
		deal.End = getMoreleEndDate(e.DOM.Find(".promo-box-countdown").Attr("data-date-to"))
		logger.Infof("Parsed end date :%s", deal.End)
	})

	c.OnRequest(func(r *colly.Request) {
		logger.Infof("Visiting %s", r.URL.String())
	})

	err := c.Visit("https://www.morele.net")
	if err != nil {
		logger.Errorf("Could not parse https://www.morele.net")
		logger.Error(err.Error())
	}
	marshall, _ := json.MarshalIndent(deal, "", "\t")
	logger.Infof("Scrapped object:\n%s", string(marshall))
	logger.Info("----------------------------------------------------------------------------------------")
	return deal
}

func getMoreleEndDate(date string, exists bool) time.Time {
	value, _ := time.Parse("2006-01-02 15:04:05", date)
	return value
}
