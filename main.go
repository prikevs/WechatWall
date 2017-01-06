package main

import (
	"crawler/config"
	"crawler/filter"
	"crawler/icrawler"
	"crawler/logger"
	"crawler/ucrawler"

	"time"
)

var log = logger.GetLogger()

func start_icrawler(usersch_filtered chan []ucrawler.User) {
	cfg := config.New()
	go icrawler.Run(cfg, usersch_filtered)
}

func start_ucrawler(usersch chan []ucrawler.User) {
	cfg := config.New()
	go func() {
		go ucrawler.Run(cfg, usersch)
		d := time.Duration(cfg.CrawlInterval) * time.Second
		for t := range time.Tick(d) {
			log.Info("user crawler starts at ", t)
			go ucrawler.Run(cfg, usersch)
		}
	}()
}

func start_filter(usersch, usersch_filtered chan []ucrawler.User) {
	cfg := config.New()
	go filter.Run(cfg, usersch, usersch_filtered)
}

func main() {
	cfg := config.New()
	usersch := make(chan []ucrawler.User, cfg.PoolSize)
	usersch_filtered := make(chan []ucrawler.User, cfg.PoolSize)

	start_icrawler(usersch_filtered)
	start_filter(usersch, usersch_filtered)
	start_ucrawler(usersch)

	forever := make(chan bool)
	<-forever
}
