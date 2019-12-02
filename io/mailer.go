package io

import (
	"deal-hunter/scrapers"
	"fmt"
	"github.com/aymerick/raymond"
	"github.com/jordan-wright/email"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/smtp"
)

const DateFormat = "2006-01-02	15:04"

func SendMail(deals []*scrapers.Deal) {
	e := email.NewEmail()
	e.From = viper.GetString("mail.source")
	e.To = viper.GetStringSlice("mail.target")
	e.Subject = "GO HotShots"
	e.HTML = getContent(deals)
	e.Send("smtp.gmail.com:587", smtp.PlainAuth("", viper.GetString("mail.login"), viper.GetString("mail.password"), "smtp.gmail.com"))
}

func getContent(deals []*scrapers.Deal) []byte {
	outerSrc, _ := ioutil.ReadFile("resources/templates/eMailMessage.html")
	src, _ := ioutil.ReadFile("resources/templates/entry.html")
	source := string(src)
	outerSource := string(outerSrc)

	ctxList := []map[string]string{}

	for _, deal := range deals {
		ctxList = append(ctxList, map[string]string{
			"price":        fmt.Sprintf("%.2f zł <s>%.2f zł</s>", deal.NewPrice, deal.OldPrice),
			"name":         deal.Name,
			"fromto":       fmt.Sprintf("%s : %s", deal.Start.Format(DateFormat), deal.End.Format(DateFormat)),
			"discountCode": deal.Code,
			"imageUrl":     deal.ImgLink,
			"itemLink":     deal.Link,
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
