package scrapers

import (
	"encoding/json"
	"github.com/google/logger"
	"reflect"
)

func logStartScrap(link string, object interface{}) {
	logger.Info("----------------------------------------------------------------------------------------")
	logger.Infof("Start parsing %s in %s", link, reflect.TypeOf(object).Name())
	logger.Info("New collector init")
}

func logDeal(deal *Deal) {
	marshall, _ := json.MarshalIndent(deal, "", "\t")
	logger.Infof("Scrapped object:\n%s", string(marshall))
	logger.Info("----------------------------------------------------------------------------------------")
}
