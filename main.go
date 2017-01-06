package main

import (
	"crawler/config"
	"crawler/filter"
	"crawler/icrawler"
	"crawler/logger"
	"crawler/ucrawler"
	"crawler/utils"

	"flag"
	"path"
	"time"
)

var log = logger.GetLogger()

type Options struct {
	CfgDir string
}

func getCurrentDir() string {
	dir, err := utils.CurrentDir()
	if err != nil {
		return ""
	}
	return dir
}

func MustParseArgs() *Options {
	cfgdir := flag.String(
		"c",
		path.Join(getCurrentDir(), "etc"),
		"directory of config files, like ./etc")

	flag.Parse()

	opts := &Options{
		CfgDir: *cfgdir,
	}
	return opts
}

func start_icrawler(cfg config.Config, usersch_filtered chan []ucrawler.User) {
	go icrawler.Run(&cfg, usersch_filtered)
}

func start_ucrawler(cfg config.Config, usersch chan []ucrawler.User) {
	go func() {
		go ucrawler.Run(&cfg, usersch)
		d := time.Duration(cfg.CrawlInterval) * time.Second
		for t := range time.Tick(d) {
			log.Info("user crawler starts at ", t)
			go ucrawler.Run(&cfg, usersch)
		}
	}()
}

func start_filter(cfg config.Config, usersch, usersch_filtered chan []ucrawler.User) {
	go filter.Run(&cfg, usersch, usersch_filtered)
}

func main() {
	opts := MustParseArgs()
	cfg := config.New(opts.CfgDir)
	usersch := make(chan []ucrawler.User, cfg.PoolSize)
	usersch_filtered := make(chan []ucrawler.User, cfg.PoolSize)

	start_icrawler(*cfg, usersch_filtered)
	start_filter(*cfg, usersch, usersch_filtered)
	start_ucrawler(*cfg, usersch)

	forever := make(chan bool)
	<-forever
}
