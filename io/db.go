package io

import (
	"deal-hunter/scrapers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/logger"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"log"
)

var db *sqlx.DB

func Subscribe(str string) bool {
	item, err := findSubscriber(db, str)
	if err != nil {
		logger.Error(err)
		return false
	}
	if item != nil && str == item.Id {
		return true
	}

	_, err = db.Exec("INSERT INTO subscribers ( subscriber_chat ) VALUES (?)", str)

	if err != nil {
		logger.Error(err)
		return false
	}
	return true
}

func AddDeal(deal *scrapers.Deal) bool {
	_, err := db.NamedExec("INSERT INTO deal ( name, link, img_link, old_price, new_price, items_sold, "+
		"items_left, start_date, end_date, promo_code ) VALUES (:name,:link, :img_link, :old_price, :new_price,"+
		" :items_sold, :items_left, :start_date, :end_date, :promo_code )", deal)

	if err != nil {
		logger.Error(err)
		return false
	}
	return true
}

func FindDeal(deal *scrapers.Deal) *scrapers.Deal {
	var deals []*scrapers.Deal
	named, err := db.PrepareNamed("SELECT * FROM deal WHERE name=:name AND link=:link AND img_link=:img_link AND old_price=:old_price AND new_price=:new_price AND end_date=:end_date AND promo_code=:promo_code ")
	if named == nil {
		logger.Error(err)
		return nil
	}
	err = named.Select(&deals, deal)

	if err != nil {
		logger.Error(err)
		return nil
	}
	if len(deals) > 0 {
		return deals[0]
	}
	return nil
}

func InitDB() {
	username := viper.GetString("db.login")
	password := viper.GetString("db.password")
	server := viper.GetString("db.server")
	name := viper.GetString("db.name")
	var err error
	db, err = sqlx.Open("mysql", username+":"+password+"@tcp("+server+":3306)/"+name+"?parseTime=true")

	if err != nil {
		logger.Error(err)
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		logger.Error(err)
		log.Panic(err)
	}
}

func findSubscriber(db *sqlx.DB, str string) (*scrapers.Subscriber, error) {
	var subscribers []*scrapers.Subscriber
	err := db.Select(&subscribers, "SELECT * FROM subscribers WHERE subscriber_chat = ?", str)
	if len(subscribers) > 0 {
		return subscribers[0], err
	}
	return nil, err
}

func FindAllSubscribers() ([]*scrapers.Subscriber, error) {
	var subscribers []*scrapers.Subscriber
	err := db.Select(&subscribers, "SELECT subscriber_chat FROM subscribers")
	if err != nil {
		panic(err.Error())
	}
	return subscribers, nil
}
