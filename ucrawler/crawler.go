package ucrawler

import (
	"crawler/config"
	"crawler/logger"
)

var log = logger.GetLogger()

func Run(cfg *config.Config, usersch chan []User) {
	fetcher := NewFetcher(cfg)

	resp, err := fetcher.Do()
	if err != nil {
		log.Error("failed to fetch, ", err.Error())
		log.Debug("config: ", *cfg)
		return
	}

	users, err := Parse(resp)
	if err != nil {
		log.Error("failed to parse, ", err.Error())
		log.Debug("resp string: ", string(resp), "\nbytes: ", resp)
		return
	}

	usersch <- users
}
