package main

import (
	"deal-hunter/io"
	"deal-hunter/scrapers"
	"flag"
	"fmt"
	"github.com/google/logger"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

var bot io.TelegramBot

var configFolder = flag.String("configFolder", ".", "full path to config folder")
var logFolder = flag.String("logFolder", "logs", "full path to log folder")

func main() {
	flag.Parse()
	defer setUpLogger().Close()

	logger.Info("Application start")
	getConfig()

	scrapers.GetDealsAt23()

	//io.InitDB()
	//
	//bot = io.TelegramBot{}
	//bot.Setup()
	//
	//logger.Info("Scheduler initialization")
	//sched := clockwork.NewScheduler()
	//logger.Info("New job every dat at 10:01 : jobAt10")
	//sched.Schedule().Every().Day().At("10:01").Do(jobAt10)
	//logger.Info("New job every dat at 22:01 : jobAt10")
	//sched.Schedule().Every().Day().At("22:01").Do(jobAt10)
	//logger.Info("New job every dat at 23:10 : jobAt23")
	//sched.Schedule().Every().Day().At("23:10").Do(jobAt23)
	//logger.Info("New job every 15 minutes : frequentJob")
	//sched.Schedule().Every(15).Minutes().Do(frequentJob)
	//logger.Info("Scheduler start: begin")
	//sched.Run()
	//logger.Info("Scheduler start: end")
}

func getConfig() {
	logger.Info("Config parsing start")

	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(*configFolder)
	viper.AddConfigPath("resources") // optionally look for config in the working directory
	err := viper.ReadInConfig()      // Find and read the config file
	if err != nil {                  // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	logger.Info("Config parsing end")
}

func jobAt10() {
	logger.Info("Start job at 10")
	handleDeals(scrapers.GetDealsAt1022())
}

func jobAt23() {
	logger.Info("Start job at 23")
	handleDeals(scrapers.GetDealsAt23())
}

func frequentJob() {
	logger.Info("Started frequent job")
	handleDeals(scrapers.GetFlexibleDeals())
}

func handleDeals(deals []*scrapers.Deal) {
	for i, deal := range deals {
		if io.FindDeal(deal) == nil {
			io.AddDeal(deal)
			bot.SendDeal(deal, i+1 == len(deals))
		}
	}
	logger.Info("End")
}

func setUpLogger() *logger.Logger {
	_ = os.Mkdir(*logFolder, os.ModeDir)
	_ = os.Chmod(*logFolder, os.ModePerm)
	filename := *logFolder + "/" + time.Now().Format("2006_01_02-15_04") + ".log"
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Error(err)
	}
	fileLogger := logger.Init(filename, true, true, f)
	logger.SetFlags(log.LstdFlags)

	return fileLogger
}
