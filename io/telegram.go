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
	logger.Infof("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	go object.handleRequests(updates)
}

func (object *TelegramBot) handleRequests(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if update.Message.Text == "/sub" || update.Message.Text == "/subscribe" {
			Subscribe(getChatID(update))
			object.reply(update, "Subscription accepted!")
		} else if update.Message.Text == "/deal" {
			object.sendDealsToSubscriber(&scrapers.Subscriber{Id: getChatID(update)}, FindLatestUniqueDeals())
		}
	}
}

func getChatID(update tgbotapi.Update) string {
	return strconv.Itoa(int(update.Message.Chat.ID))
}

func (object *TelegramBot) reply(update tgbotapi.Update, message string) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	msg.ReplyToMessageID = update.Message.MessageID
	msg.DisableWebPagePreview = true

	_, _ = object.bot.Send(msg)
}

func (object *TelegramBot) SendDealsToAll(deals []*scrapers.Deal) {
	subscribers, err := FindAllSubscribers()
	if err != nil {
		logger.Error(err)
		return
	}
	for _, subscriber := range subscribers {
		object.sendDealsToSubscriber(subscriber, deals)
	}
}

func (object *TelegramBot) sendDealsToSubscriber(subscriber *scrapers.Subscriber, deals []*scrapers.Deal) {
	strint, _ := strconv.Atoi(subscriber.Id)
	for _, deal := range deals {
		msg := tgbotapi.NewMessage(int64(strint), getDealMessage(deal))
		msg.DisableWebPagePreview = true
		_, err := object.bot.Send(msg)
		object.handleMsgSent(err, deal, subscriber)
	}
	object.sendEndMsg(int64(strint))
}

func (object *TelegramBot) handleMsgSent(err error, deal *scrapers.Deal, subscriber *scrapers.Subscriber) {
	if err != nil {
		logger.Error("ERROR WHILE SENDING DEAL : " + deal.Name)
		logger.Error(err)
		logger.Error(err.Error())
	} else {
		logger.Info("DEAL SENT : " + subscriber.Id + " : " + deal.Name)
	}
}

func (object *TelegramBot) SendDeal(deal *scrapers.Deal, withEnding bool) {
	subscribers, err := FindAllSubscribers()
	if err != nil {
		logger.Error(err)
		return
	}
	for _, subscriber := range subscribers {
		strint, _ := strconv.Atoi(subscriber.Id)
		msg := tgbotapi.NewMessage(int64(strint), getDealMessage(deal))
		msg.DisableWebPagePreview = true
		_, err := object.bot.Send(msg)
		object.handleMsgSent(err, deal, subscriber)
		if withEnding {
			object.sendEndMsg(int64(strint))
		}
	}
}

func (object *TelegramBot) sendEndMsg(chatId int64) {
	msg := tgbotapi.NewMessage(chatId, "\xE2\x9C\x85 That's all I have for you, my master! \xE2\x9C\x85")
	msg.DisableWebPagePreview = true
	_, _ = object.bot.Send(msg)
}

func getDealMessage(deal *scrapers.Deal) string {
	return "name : " + deal.Name +
		"\nprice:" + fmt.Sprintf("    %.2f zł // %.2f zł", deal.NewPrice, deal.OldPrice) +
		"\ndiscountCode: " + deal.Code +
		"\nitemLink: " + deal.Link
}
