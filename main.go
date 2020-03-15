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
	getConfig()

	io.InitDB()

	bot = io.TelegramBot{}
	bot.Setup()

	logger.Info("Scheduler initialization")
	sched := clockwork.NewScheduler()
	logger.Info("New job every dat at 10:01 : jobAt10")
	sched.Schedule().Every().Day().At("10:01").Do(jobAt10)
	logger.Info("New job every dat at 22:01 : jobAt10")
	sched.Schedule().Every().Day().At("22:01").Do(jobAt10)
	logger.Info("New job every dat at 23:10 : jobAt23")
	sched.Schedule().Every().Day().At("23:10").Do(jobAt23)
	logger.Info("Scheduler start: begin")
	sched.Run()
	logger.Info("Scheduler start: end")
}

func getConfig() {
	logger.Info("Config parsing start")

	viper.SetConfigName("config")    // name of config file (without extension)
	viper.AddConfigPath(".")         // optionally look for config in the working directory
	viper.AddConfigPath("resources") // optionally look for config in the working directory
	err := viper.ReadInConfig()      // Find and read the config file
	if err != nil {                  // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	logger.Info("Config parsing end")
}

func jobAt10() {
	logger.Info("Start job at 10")
	deals := scrapers.GetDealsAt1022()
	for i, deal := range deals {
		if io.FindDeal(deal) == nil {
			io.AddDeal(deal)
			bot.SendDeal(deal, i+1 == len(deals))
		}
	}
	logger.Info("End")
}

func jobAt23() {
	logger.Info("Start job at 23")
	deals := scrapers.GetDealsAt23()
	bot.SendDeals(deals)
	for _, deal := range deals {
		io.AddDeal(deal)
	}
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
