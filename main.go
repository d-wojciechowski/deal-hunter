package main

import (
	"deal-hunter/io"
	"deal-hunter/scrapers"
	"fmt"
	"github.com/spf13/viper"
)

func main() {
	getConfig()
	kom := scrapers.ScrapXKom()
	io.SendMail([]*scrapers.Deal{kom})
}

func getConfig() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
