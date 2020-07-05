package scrapers

import (
	"encoding/json"
	"github.com/google/logger"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type CombatScrapper struct {
	URL *url.URL
}

func (scrapper *CombatScrapper) Scrap() *Deal {
	logStartScrap(scrapper.URL.String(), scrapper)
	deal := &Deal{SiteName: "combat"}

	err, rawJson := scrapper.getJSON()
	if err != nil {
		logger.Error("Could not deserialize json. Root error: ")
		logger.Error(err.Error())
	} else {
		json.Unmarshal(rawJson["name"], &deal.Name)
		json.Unmarshal(rawJson["regular_url"], &deal.Link)
		json.Unmarshal(rawJson["photo"], &deal.ImgLink)
		deal.ImgLink = scrapper.URL.String() + "/pub/media/catalog/product" + deal.ImgLink

		tempString := ""
		json.Unmarshal(rawJson["regular_price"], &tempString)
		tempString = tempString[strings.Index(tempString, ">")+1 : strings.Index(tempString, "\u00a0z")]
		deal.OldPrice, _ = convertToNumber(tempString)

		json.Unmarshal(rawJson["promotion_price"], &tempString)
		tempString = tempString[strings.Index(tempString, ">")+1 : strings.Index(tempString, "\u00a0z")]
		deal.NewPrice, _ = convertToNumber(tempString)

		now := time.Now()
		deal.Start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		deal.End = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

		json.Unmarshal(rawJson["left"], &deal.Left)
		json.Unmarshal(rawJson["sold"], &deal.Sold)

	}

	logDeal(deal)
	return deal
}

func (scrapper *CombatScrapper) getJSON() (err error, rawJson map[string]json.RawMessage) {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", scrapper.URL.String()+"/rest/V1/get-hot-shot", nil)
	request.Header.Add("Content-Type", "application/json")
	response, _ := client.Do(request)
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body[1:len(body)-1], &rawJson)
	return
}

func convertToNumber(value string) (float64, error) {
	price := MoneyRegexp.ReplaceAllString(value, "")
	price = strings.Replace(price, ",", ".", -1)
	return strconv.ParseFloat(price, 64)
}
