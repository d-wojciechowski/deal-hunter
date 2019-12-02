package main

import (
	"deal-hunter/io"
	"deal-hunter/scrapers"
	"fmt"
	"github.com/spf13/viper"
	"github.com/whiteShtef/clockwork"
)

func main() {
	getConfig()

	jobAt10()

	sched := clockwork.NewScheduler()
	sched.Schedule().Every().Day().At("10:01").Do(jobAt10)
	sched.Schedule().Every().Day().At("22:01").Do(jobAt10)
	sched.Run()
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
	kom := scrapers.ScrapXKomGroup("https://www.x-kom.pl/")
	alto := scrapers.ScrapXKomGroup("https://www.al.to/")
	morele := scrapers.ScrapMorele()
	io.SendMail([]*scrapers.Deal{kom, alto, morele})
}
