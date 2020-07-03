package scheduler

import (
	"deal-hunter/io"
	"deal-hunter/scrapers"
	"github.com/google/logger"
	"github.com/spf13/viper"
	"github.com/whiteShtef/clockwork"
	"strconv"
	"strings"
)

var jobConfig = make(map[string]func() *scrapers.Deal)
var bot = io.TelegramBot{}

func InitJobs() {
	jobConfig["xkom"] = func() *scrapers.Deal {
		return scrapers.ScrapXKomGroup("https://www.x-kom.pl/")
	}
	jobConfig["alto"] = func() *scrapers.Deal {
		return scrapers.ScrapXKomGroup("https://www.al.to/")
	}
	jobConfig["combat"] = scrapers.ScrapCombat
	jobConfig["morele"] = scrapers.ScrapMorele

	bot.Setup()
}

type ScheduleJob struct {
	Handler func([]*scrapers.Deal)
	Steps   []func() *scrapers.Deal
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
			Steps:   make([]func() *scrapers.Deal, 0),
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
		stepResult := step()
		deals = append(deals, stepResult)
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
