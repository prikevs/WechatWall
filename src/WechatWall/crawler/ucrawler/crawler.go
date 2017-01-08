package ucrawler

import (
	"WechatWall/crawler/config"
	"WechatWall/logger"
)

var log = logger.GetLogger("crawler/ucrawler")

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
		log.Error("failed to parse,", err.Error(), string(resp))
		log.Debug("resp bytes:", resp)
		return
	}

	usersch <- users
}
