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
	logger.Info("----------------------------------------------------------------------------------------")
	logger.Infof("Start parsing %s in ScrapCombat", scrapper.URL.String())
	deal := &Deal{}

	client := &http.Client{}
	request, _ := http.NewRequest("GET", scrapper.URL.String()+"/rest/V1/get-hot-shot", nil)
	request.Header.Add("Content-Type", "application/json")
	response, _ := client.Do(request)
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	var objmap map[string]json.RawMessage
	err := json.Unmarshal(body[1:len(body)-1], &objmap)
	if err != nil {
		logger.Error("Could not deserialize json. Root error: ")
		logger.Error(err.Error())
	}

	deal.SiteName = "combat"
	json.Unmarshal(objmap["name"], &deal.Name)
	json.Unmarshal(objmap["regular_url"], &deal.Link)
	json.Unmarshal(objmap["photo"], &deal.ImgLink)
	deal.ImgLink = scrapper.URL.String() + "/pub/media/catalog/product" + deal.ImgLink

	tempString := ""
	json.Unmarshal(objmap["regular_price"], &tempString)
	tempString = tempString[strings.Index(tempString, ">")+1 : strings.Index(tempString, "\u00a0z")]
	deal.OldPrice, _ = convertToNumber(tempString)

	json.Unmarshal(objmap["promotion_price"], &tempString)
	tempString = tempString[strings.Index(tempString, ">")+1 : strings.Index(tempString, "\u00a0z")]
	deal.NewPrice, _ = convertToNumber(tempString)

	now := time.Now()
	deal.Start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	deal.End = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	json.Unmarshal(objmap["left"], &deal.Left)
	json.Unmarshal(objmap["sold"], &deal.Sold)

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
