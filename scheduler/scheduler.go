package scheduler

import (
	"deal-hunter/io"
	"deal-hunter/scrapers"
	"github.com/google/logger"
	"github.com/spf13/viper"
	"github.com/whiteShtef/clockwork"
	"net/url"
	"strconv"
	"strings"
)

var jobConfig = make(map[string]scrapers.Scrapper)
var bot = io.TelegramBot{}

func InitJobs() {
	jobConfig["xkom"] = &scrapers.XKomGroupScrapper{URL: parseURL("https://x-kom.pl")}
	jobConfig["alto"] = &scrapers.XKomGroupScrapper{URL: parseURL("https://al.to")}
	jobConfig["combat"] = &scrapers.CombatScrapper{URL: parseURL("https://combat.pl")}
	jobConfig["morele"] = &scrapers.MoreleScrapper{URL: parseURL("https://morele.net")}

	bot.Setup()
}

func parseURL(link string) *url.URL {
	parse, err := url.Parse(link)
	if err != nil {
		logger.Errorf("Could not parse given URL: %s", link)
		panic(err)
	}
	return parse
}

type ScheduleJob struct {
	Handler func([]*scrapers.Deal)
	Steps   []scrapers.Scrapper
}

func CreateScheduler() *clockwork.Scheduler {
	dailyJSON := viper.GetStringMap("scheduler.daily")
	sched := clockwork.NewScheduler()
	constructJob(dailyJSON, func(interval string, executionJob func()) {
		sched.Schedule().Every().Day().At(interval).Do(executionJob)
	})
	everyJson := viper.GetStringMap("scheduler.every")
	constructJob(everyJson, func(interval string, executionJob func()) {
		every, _ := strconv.Atoi(interval)
		sched.Schedule().Every(every).Minutes().Do(executionJob)
	})

	return &sched
}

func constructJob(stringMap map[string]interface{}, function func(interval string, executionJob func())) {
	for key, value := range stringMap {
		job := ScheduleJob{
			Handler: defaultHandler,
			Steps:   make([]scrapers.Scrapper, 0),
		}
		for _, entry := range value.([]interface{}) {
			job.Steps = append(job.Steps, jobConfig[entry.(string)])
		}
		function(strings.Replace(key, "-", ":", -1), job.execute)
	}
}

func (job *ScheduleJob) execute() {
	deals := make([]*scrapers.Deal, 0)
	for _, step := range job.Steps {
		deals = append(deals, step.Scrap())
	}
	job.Handler(deals)
}

func defaultHandler(deals []*scrapers.Deal) {
	for i, deal := range deals {
		if deal != nil && deal.Name != "" && io.FindDeal(deal) == nil {
			io.AddDeal(deal)
			bot.SendDeal(deal, i+1 == len(deals))
		}
	}
	logger.Info("End")
}
