package icrawler

import (
	"crawler/config"
	"crawler/ucrawler"
	"crawler/utils"

	"log"
	// "time"
)

func Download(cfg *config.Config, user *ucrawler.User) {
	ifetcher := NewIFetcher(cfg, user)
	resp, err := ifetcher.Do()
	if err != nil {
		// Log
		log.Println("Download failed, " + err.Error())
		return
	}
	dir := utils.BuildImagePath(cfg, user)
	if err := utils.WriteFile(dir, resp); err != nil {
		log.Println("Download failed, " + err.Error())
		// Log
	}
}

func Worker(wid int, cfg *config.Config, userch <-chan ucrawler.User, exited chan<- int, exit <-chan bool) {
	// Log: start one worker
	log.Println("start one worker ", wid)
	for {
		select {
		case user := <-userch:
			log.Printf("worker %d start to download %s, %s",
				wid, user.UserName, user.UserOpenid)
			Download(cfg, &user)
		case <-exit:
			log.Println("worker ", wid, " received signal to stop")
			exited <- wid
			return
		}
	}
}

func RunPool(cfg *config.Config, usersch chan []ucrawler.User,
	runable func(int, *config.Config, <-chan ucrawler.User, chan<- int, <-chan bool)) {

	// start workers
	log.Println("start workers")

	userch := make(chan ucrawler.User, cfg.PoolSize)
	exited := make(chan int, cfg.PoolSize)
	exit := make(chan bool, cfg.PoolSize)
	for i := 1; i <= cfg.PoolSize; i++ {
		go Worker(i, cfg, userch, exited, exit)
	}

	// handle list from ucrawler.User
	for users := range usersch {
		if len(users) == 0 {
			log.Println("worker master received signal to stop")
			// send nil to all workers
			for i := 0; i < cfg.PoolSize; i++ {
				exit <- true
			}

			// Log: wait all workers to stop
			log.Println("wait all workers to stop")
			count := 0
			for wid := range exited {
				log.Printf("worker %d stopped.", wid)
				count++
				if count == cfg.PoolSize {
					log.Println("all workers stopped, exiting pool")
					return
				}
			}
		}
		log.Println(len(users))
		for _, user := range users {
			// log.Println(users[i])
			userch <- user
			// time.Sleep(time.Second)
		}
	}
}

func Run(cfg *config.Config, usersch chan []ucrawler.User) {
	RunPool(cfg, usersch, Worker)
}
