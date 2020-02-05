package io

import (
	"deal-hunter/scrapers"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/logger"
	"github.com/spf13/viper"
	"log"
	"strconv"
)

type TelegramBot struct {
	bot *tgbotapi.BotAPI
}

func (object *TelegramBot) Setup() {
	logger.Info("Telegram setup start")
	token := viper.GetString("telegram.token")
	bot, err := tgbotapi.NewBotAPI(token)
	object.bot = bot
	if err != nil {
		log.Panic(err)
	}

	logger.Info("Telegram setup completed")
	logger.Info("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	go func() {
		for update := range updates {
			if update.Message == nil { // ignore any non-Message Updates
				continue
			}

			if update.Message.Text == "/sub" || update.Message.Text == "/subscribe" {
				Subscribe(strconv.Itoa(int(update.Message.Chat.ID)))

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "SUBSCRIBED!")
				msg.ReplyToMessageID = update.Message.MessageID

				object.bot.Send(msg)
			}
		}
	}()
}

func (object *TelegramBot) SendDeal(deals []*scrapers.Deal) {
	subscribers, err := FindAllSubscribers()
	if err != nil {
		logger.Error(err)
		return
	}
	for _, deal := range deals {
		for _, subscriber := range subscribers {
			strint, _ := strconv.Atoi(subscriber)
			msgTxt := "price:" + fmt.Sprintf("    %.2f zł // %.2f zł", deal.NewPrice, deal.OldPrice) +
				"\ndiscountCode: " + deal.Code +
				"\nitemLink: " + deal.Link
			msg := tgbotapi.NewMessage(int64(strint), msgTxt)
			object.bot.Send(msg)
			logger.Info("DEAL SENT : " + subscriber + " : " + deal.Name)
		}
	}
}
