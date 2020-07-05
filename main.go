package main

import (
	"deal-hunter/io"
	"deal-hunter/scheduler"
	"flag"
	"fmt"
	"github.com/google/logger"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

var configFolder = flag.String("configFolder", ".", "full path to config folder")
var logFolder = flag.String("logFolder", "logs", "full path to log folder")

func main() {
	flag.Parse()
	defer setUpLogger().Close()

	logger.Info("Application start")
	getConfig()

	io.InitDB()

	scheduler.InitJobs()
	scheduler.CreateScheduler().Run()
	logger.Info("Scheduler start: end")
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
