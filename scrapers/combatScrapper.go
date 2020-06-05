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
	baseSite := "https://www.combat.pl/"

	c.OnHTML(".main", func(e *colly.HTMLElement) {
		deal.SiteName = "combat"
		deal.Name = strings.Replace(e.DOM.Find(".page-title > span").Text(), "\n", "", -1)
		logger.Infof("Parsed name :%s", deal.Name)
		s2 := e.DOM.Find(".sku > span").Last()
		deal.Link = baseSite + s2.Text()
		logger.Infof("Parsed link :%s", deal.Link)
		s3 := e.DOM.Find("img")
		deal.ImgLink, _ = s3.Attr("src")
		deal.ImgLink = strings.Replace(deal.ImgLink, "?filters=grayscale", "", -1)
		logger.Infof("Parsed img link :%s", deal.ImgLink)
		find := e.DOM.Find(".price")
		deal.OldPrice, _ = convertToNumber(find.Nodes[1].FirstChild.Data)
		logger.Infof("Parsed old price :%0.2f", deal.OldPrice)
		deal.NewPrice, _ = convertToNumber(find.Nodes[0].FirstChild.Data)
		logger.Infof("Parsed new price :%0.2f", deal.NewPrice)

		deal.Start = getStartDate().Add(time.Duration(1) * time.Hour)
		logger.Infof("Parsed start date :%s", deal.Start)
		deal.End = getStartDate().Add(time.Duration(25) * time.Hour)
		logger.Infof("Parsed end date :%s", deal.End)
	})

	c.OnRequest(func(r *colly.Request) {
		logger.Infof("Visiting %s", r.URL.String())
	})

	err := c.Visit(baseSite + "goracy-strzal")
	if err != nil {
		logger.Errorf("Could not parse %sgoracy-strzal", baseSite)
		logger.Error(err.Error())
	}
	marshall, _ := json.MarshalIndent(deal, "", "\t")
	logger.Infof("Scrapped object:\n%s", string(marshall))
	logger.Info("----------------------------------------------------------------------------------------")
	return deal
}

func convertToNumber(value string) (float64, error) {
	price := MoneyRegexp.ReplaceAllString(value, "")
	price = strings.Replace(price, ",", ".", -1)
	return strconv.ParseFloat(price, 64)
}
