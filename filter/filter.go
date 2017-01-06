package filter

import (
	"crawler/config"
	"crawler/logger"
	"crawler/ucrawler"
	"crawler/utils"
)

var log = logger.GetLogger()

func need(cfg *config.Config, user *ucrawler.User) bool {
	path := utils.BuildImagePath(cfg, user)
	if utils.FileExists(path) {
		return false
	}
	return true
}

func Run(cfg *config.Config, usersch <-chan []ucrawler.User, usersch_filtered chan<- []ucrawler.User) {
	for users := range usersch {
		filtered := make([]ucrawler.User, 0)
		for _, user := range users {
			if need(cfg, &user) {
				filtered = append(filtered, user)
			}
		}
		num := len(filtered)
		if num == 0 {
			log.Info("all users have been fetched, nothing to do")
			continue
		} else {
			log.Info(num, " more images need to be crawled")
		}
		usersch_filtered <- filtered
	}
}
