package main

import (
	"deal-hunter/io"
	"deal-hunter/scrapers"
	"fmt"
	"github.com/google/logger"
	"github.com/spf13/viper"
	"github.com/whiteShtef/clockwork"
	"log"
	"os"
	"time"
)

var bot io.TelegramBot

func main() {
	defer setUpLogger().Close()

	logger.Info("Application start")
	logger.Info("Config parsing start")
	getConfig()
	logger.Info("Config parsing end")
	bot = io.TelegramBot{}
	bot.Setup()
	jobAt10()

	logger.Info("Scheduler initialization")
	sched := clockwork.NewScheduler()
	logger.Info("New job every dat at 10:01 : jobAt10")
	sched.Schedule().Every().Day().At("10:01").Do(jobAt10)
	logger.Info("New job every dat at 22:01 : jobAt10")
	sched.Schedule().Every().Day().At("22:01").Do(jobAt10)
	logger.Info("Scheduler start: begin")
	sched.Run()
	logger.Info("Scheduler start: end")
}

func getConfig() {
	viper.SetConfigName("config")    // name of config file (without extension)
	viper.AddConfigPath(".")         // optionally look for config in the working directory
	viper.AddConfigPath("resources") // optionally look for config in the working directory
	err := viper.ReadInConfig()      // Find and read the config file
	if err != nil {                  // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func jobAt10() {
	logger.Info("Start job at 10")
	kom := scrapers.ScrapXKomGroup("https://www.x-kom.pl/")
	alto := scrapers.ScrapXKomGroup("https://www.al.to/")
	morele := scrapers.ScrapMorele()
	io.SendMail([]*scrapers.Deal{kom, alto, morele})
	bot.SendDeal([]*scrapers.Deal{kom, alto, morele})
	logger.Info("End")
}

func setUpLogger() *logger.Logger {
	_ = os.Mkdir("logs", os.ModeDir)
	_ = os.Chmod("logs", os.ModePerm)
	filename := "logs/" + time.Now().Format("2006_01_02-15_04") + ".log"
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Error(err)
	}
	fileLogger := logger.Init(filename, true, true, f)
	logger.SetFlags(log.LstdFlags)

	return fileLogger
}
