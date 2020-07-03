package scrapers

import (
	"encoding/json"
	"github.com/gocolly/colly"
	"github.com/google/logger"
	"regexp"
	"strings"
	"time"
)

var MoneyRegexp = regexp.MustCompile("[  zł]")

func ScrapXKomGroup(root string) *Deal {
	logger.Info("----------------------------------------------------------------------------------------")
	logger.Infof("Start parsing %s in XKomGroupParser", root)
	logger.Info("New collector init")
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.97 Safari/537.36 OPR/65.0.3467.48"),
	)
	deal := &Deal{}
	hotShotRoot := root + "/goracy_strzal"

	c.OnHTML("html", func(e *colly.HTMLElement) {
		if strings.Contains(hotShotRoot, "x-kom") {
			deal.SiteName = "x-kom"
		} else {
			deal.SiteName = "al.to"
		}

		scripts := e.DOM.Find("script")
		var jsonData string
		for _, script := range scripts.Nodes {
			child := script.FirstChild
			if child != nil {
				data := child.Data
				if strings.Contains(data, "hotShot") {
					jsonData = data
					break
				}
			}
		}

		jsonData = jsonData[strings.Index(jsonData, "hotShot") : len(jsonData)-1]
		jsonData = jsonData[strings.Index(jsonData, "extend")+8 : len(jsonData)-1]
		jsonData = jsonData[0 : strings.Index(jsonData, "},\"id\"")+1]

		var objmap map[string]json.RawMessage
		err := json.Unmarshal([]byte(jsonData), &objmap)
		if err != nil {
			panic(err)
		}
		json.Unmarshal(objmap["promotionName"], &deal.Name)
		logger.Infof("Parsed name :%s", deal.Name)

		json.Unmarshal(objmap["price"], &deal.NewPrice)
		logger.Infof("Parsed new price :%0.2f", deal.NewPrice)

		json.Unmarshal(objmap["oldPrice"], &deal.OldPrice)
		logger.Infof("Parsed old price :%0.2f", deal.OldPrice)

		dealUrl, _ := e.DOM.Find("meta[property='og:url']").Attr("content")
		deal.Link = dealUrl
		logger.Infof("Parsed link :%s", deal.Link)

		dealImageLink, _ := e.DOM.Find("meta[property='og:image']").Attr("content")
		deal.ImgLink = dealImageLink
		deal.ImgLink = strings.Replace(deal.ImgLink, "?filters=grayscale", "", -1)
		logger.Infof("Parsed img link :%s", deal.ImgLink)

		json.Unmarshal(objmap["promotionTotalCount"], &deal.Left)
		logger.Infof("Parsed left count :%d", deal.Left)

		e.DOM.Find("script")
		deal.Start = getStartDate()
		logger.Infof("Parsed start date :%s", deal.Start)
		deal.End = getEndDate()
		logger.Infof("Parsed end date :%s", deal.End)
	})

	c.OnRequest(func(r *colly.Request) {
		logger.Infof("Visiting %s", r.URL.String())
	})

	err := c.Visit(hotShotRoot)
	if err != nil {
		logger.Errorf("Could not parse %s", hotShotRoot)
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
