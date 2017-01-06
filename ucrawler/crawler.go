package ucrawler

import (
	"crawler/config"
)

func Run(cfg *config.Config, usersch chan []User) {
	fetcher := NewFetcher(cfg)

	resp, err := fetcher.Do()
	if err != nil {
		// Log
		return
	}

	users, err := Parse(resp)
	if err != nil {
		// Log
		return
	}

	usersch <- users
}
