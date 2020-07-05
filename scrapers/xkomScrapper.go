package scrapers

import (
	"encoding/json"
	"github.com/gocolly/colly"
	"github.com/google/logger"
	"net/url"
	"strings"
	"time"
)

type XKomGroupScrapper struct {
	URL *url.URL
}

func (scrapper *XKomGroupScrapper) Scrap() *Deal {
	targetLink := scrapper.URL.String() + "/goracy_strzal"
	logStartScrap(targetLink, scrapper)

	wrapper := collyWrapper{
		Link:     targetLink,
		Selector: "html",
		hostname: scrapper.URL.Hostname(),
		Consumer: scrapper.scrap,
	}
	wrapper.init()
	deal := wrapper.execute()

	logDeal(deal)
	return deal
}

func (scrapper *XKomGroupScrapper) scrap(deal *Deal, e *colly.HTMLElement) {
	objmap := scrapper.getHotShotJSON(e)

	json.Unmarshal(objmap["promotionName"], &deal.Name)
	json.Unmarshal(objmap["price"], &deal.NewPrice)
	json.Unmarshal(objmap["oldPrice"], &deal.OldPrice)
	json.Unmarshal(objmap["promotionTotalCount"], &deal.Left)

	deal.Link, _ = e.DOM.Find("meta[property='og:url']").Attr("content")
	deal.ImgLink, _ = e.DOM.Find("meta[property='og:image']").Attr("content")
	deal.ImgLink = strings.Replace(deal.ImgLink, "?filters=grayscale", "", -1)

	deal.Start = scrapper.getStartDate()
	deal.End = scrapper.getEndDate()
}

func (scrapper *XKomGroupScrapper) getHotShotJSON(e *colly.HTMLElement) map[string]json.RawMessage {
	jsonData := ""
	for {
		jsonData = scrapper.findDirtyHotShotJSON(e)
		if jsonData == "" {
			time.Sleep(1 * time.Minute)
		} else {
			break
		}
	}

	jsonData = jsonData[strings.Index(jsonData, "hotShot") : len(jsonData)-1]
	jsonData = jsonData[strings.Index(jsonData, "extend")+8 : len(jsonData)-1]
	jsonData = jsonData[0 : strings.Index(jsonData, "},\"id\"")+1]

	var objmap map[string]json.RawMessage
	err := json.Unmarshal([]byte(jsonData), &objmap)
	if err != nil {
		logger.Error("Could not deserialize json. Root error: ")
		logger.Error(err.Error())
	}
	return objmap
}

func (scrapper *XKomGroupScrapper) findDirtyHotShotJSON(e *colly.HTMLElement) string {
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
	return jsonData
}

func (scrapper *XKomGroupScrapper) getStartDate() time.Time {
	now := time.Now()
	if now.Hour() > 9 && now.Hour() < 22 {
		return time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, now.Location())
	}
	return time.Date(now.Year(), now.Month(), now.Day(), 22, 0, 0, 0, now.Location())
}

func (scrapper *XKomGroupScrapper) getEndDate() time.Time {
	now := time.Now()
	if now.Hour() > 9 && now.Hour() < 22 {
		return time.Date(now.Year(), now.Month(), now.Day(), 21, 59, 59, 0, now.Location())
	}
	return time.Date(now.Year(), now.Month(), now.Day()+1, 9, 59, 59, 0, now.Location())
}
