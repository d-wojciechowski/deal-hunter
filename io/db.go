package io

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/logger"
	"github.com/spf13/viper"
)

func Subscribe(str string) bool {
	db, err := openConnection()
	if err != nil {
		logger.Error(err)
		return false
	}
	defer db.Close()

	item, err := findSubscriber(db, str)
	if err != nil {
		logger.Error(err)
		return false
	}
	if str == item {
		return true
	}

	insert, err := db.Query("INSERT INTO subscribers VALUES (" + str + ")")

	if err != nil {
		logger.Error(err)
		return false
	}
	defer insert.Close()
	return true
}

func openConnection() (*sql.DB, error) {
	username := viper.GetString("db.login")
	password := viper.GetString("db.password")
	server := viper.GetString("db.server")
	name := viper.GetString("db.name")

	db, err := sql.Open("mysql", username+":"+password+"@tcp("+server+":3306)/"+name)
	return db, err
}

func findSubscriber(db *sql.DB, str string) (string, error) {
	query, _ := db.Query("SELECT s.subscriber_chat FROM subscribers s WHERE s.subscriber_chat = " + str)

	var result string
	query.Next()
	err := query.Scan(&result)

	return result, err
}

func FindAllSubscribers() ([]string, error) {
	db, err := openConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	results, err := db.Query("SELECT subscriber_chat FROM subscribers")
	if err != nil {
		panic(err.Error())
	}
	strings := make([]string, 1)
	for results.Next() {
		var result string
		err := results.Scan(&result)
		if err != nil {
			return nil, err
		}
		strings = append(strings, result)
	}
	defer results.Close()
	return strings, nil
}
