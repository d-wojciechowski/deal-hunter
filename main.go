package main

import (
	"fmt"
	"github.com/aymerick/raymond"
	"github.com/gocolly/colly"
	"github.com/jordan-wright/email"
	"github.com/spf13/viper"
	"net/smtp"

	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

type Deal struct {
	Name     string
	ImgLink  string
	OldPrice float64
	NewPrice float64
	Sold     int64
	Left     int64
}

func main() {
	getConfig()

	kom := scrapXKom()
	sendMail([]*Deal{kom})
}

var moneyRegexp = regexp.MustCompile("[ zł]")

func sendMail(deals []*Deal) {
	e := email.NewEmail()
	e.From = viper.GetString("mail.source")
	e.To = viper.GetStringSlice("mail.target")
	e.Subject = "GO HotShots"
	e.HTML = getContent(deals)
	e.Send("smtp.gmail.com:587", smtp.PlainAuth("", viper.GetString("mail.login"), viper.GetString("mail.password"), "smtp.gmail.com"))
}

func getContent(deals []*Deal) []byte {
	outerSrc, _ := ioutil.ReadFile("templates/eMailMessage.html")
	src, _ := ioutil.ReadFile("templates/entry.html")
	source := string(src)
	outerSource := string(outerSrc)

	ctxList := []map[string]string{}

	for _, deal := range deals {
		ctxList = append(ctxList, map[string]string{
			"price":        fmt.Sprintf("%.2f zł <s>%.2f zł</s>", deal.NewPrice, deal.OldPrice),
			"name":         deal.Name,
			"fromto":       "date range",
			"discountCode": "123",
			"imageUrl":     deal.ImgLink,
			"itemLink":     deal.ImgLink,
		})
	}

	ctxListOuter := []map[string]string{
		{
			"content": parseHtml(source, ctxList),
		},
	}

	return []byte(parseHtml(outerSource, ctxListOuter))
}

func parseHtml(source string, ctxList []map[string]string) string {
	// parse template
	tpl, err := raymond.Parse(source)
	if err != nil {
		panic(err)
	}
	result := ""
	for _, ctx := range ctxList {
		// render template
		res, err := tpl.Exec(ctx)
		if err != nil {
			panic(err)
		}
		result += res
	}
	return result
}

func scrapXKom() *Deal {
	// Instantiate default collector
	c := colly.NewCollector()
	deal := &Deal{}

	c.OnHTML(".hot-shot", func(e *colly.HTMLElement) {
		pImpression := e.DOM.Find(".product-impression")
		deal.Name = pImpression.Find(".product-name").Text()
		deal.ImgLink, _ = pImpression.Find("img").Attr("src")
		priceDiv := e.DOM.Find(".price")

		oldPrice := moneyRegexp.ReplaceAllString(priceDiv.Find(".old-price").Text(), "")
		oldPrice = strings.Replace(oldPrice, ",", ".", 1)
		deal.OldPrice, _ = strconv.ParseFloat(oldPrice, 64)

		newPrice := moneyRegexp.ReplaceAllString(priceDiv.Find(".new-price").Text(), "")
		newPrice = strings.Replace(newPrice, ",", ".", 1)
		deal.NewPrice, _ = strconv.ParseFloat(newPrice, 64)

		countDiv := e.DOM.Find(".count")
		deal.Left, _ = strconv.ParseInt(countDiv.Find(".pull-Left > .gs-quantity").Text(), 10, 64)
		deal.Sold, _ = strconv.ParseInt(countDiv.Find(".pull-right > .gs-quantity").Text(), 10, 64)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit("https://x-kom.pl/")
	return deal
}

func getConfig() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
