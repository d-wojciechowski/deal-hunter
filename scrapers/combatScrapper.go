package scrapers

import (
	"encoding/json"
	"github.com/gocolly/colly"
	"github.com/google/logger"
	"strconv"
	"strings"
	"time"
)

func ScrapCombat() *Deal {
	logger.Info("----------------------------------------------------------------------------------------")
	logger.Infof("Start parsing https://www.combat.pl/ in ScrapCombat")
	logger.Info("New collector init")
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.97 Safari/537.36 OPR/65.0.3467.48"),
	)
	deal := &Deal{}

	c.OnHTML(".hot-shot", func(e *colly.HTMLElement) {
		deal.SiteName = "combat"
		deal.Name = strings.Replace(e.DOM.Find(".product-name").Text(), "\n", "", -1)
		logger.Infof("Parsed name :%s", deal.Name)
		deal.Link, _ = e.DOM.Find(".product-item-link").Attr("href")
		logger.Infof("Parsed link :%s", deal.Link)
		deal.ImgLink, _ = e.DOM.Find(".product-image-photo").Attr("src")
		deal.ImgLink = strings.Replace(deal.ImgLink, "?filters=grayscale", "", -1)
		logger.Infof("Parsed img link :%s", deal.ImgLink)
		priceDiv := e.DOM.Find(".price-box")

		price := MoneyRegexp.ReplaceAllString(priceDiv.Find("span > .price").Text(), "")
		price = strings.Replace(price, ",", ".", -1)
		prices := strings.Split(price, " ")
		deal.OldPrice, _ = strconv.ParseFloat(prices[1], 64)
		logger.Infof("Parsed old price :%0.2f", deal.OldPrice)

		deal.NewPrice, _ = strconv.ParseFloat(prices[0], 64)
		logger.Infof("Parsed new price :%0.2f", deal.NewPrice)

		countDiv := e.DOM.Find(".deal-stock-label")
		deal.Left, _ = strconv.ParseInt(countDiv.Find(".stock-available > strong").Text(), 10, 64)
		logger.Infof("Parsed left count :%d", deal.Left)
		deal.Sold, _ = strconv.ParseInt(countDiv.Find(".stock-sold > strong").Text(), 10, 64)
		logger.Infof("Parsed sold count :%d", deal.Sold)

		deal.Start = getStartDate().Add(time.Duration(1) * time.Hour)
		logger.Infof("Parsed start date :%s", deal.Start)
		deal.End = getStartDate().Add(time.Duration(25) * time.Hour)
		logger.Infof("Parsed end date :%s", deal.End)
	})

	c.OnRequest(func(r *colly.Request) {
		logger.Infof("Visiting %s", r.URL.String())
	})

	err := c.Visit("https://www.combat.pl/")
	if err != nil {
		logger.Error("Could not parse https://www.combat.pl/")
		logger.Error(err.Error())
	}
	marshall, _ := json.MarshalIndent(deal, "", "\t")
	logger.Infof("Scrapped object:\n%s", string(marshall))
	logger.Info("----------------------------------------------------------------------------------------")
	return deal
}
