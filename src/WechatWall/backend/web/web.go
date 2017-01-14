package web

import (
	"WechatWall/backend/config"
	"WechatWall/logger"

	//	"github.com/gorilla/mux"

	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"text/template"
)

var log = logger.GetLogger("backend/web")

var (
	imageDir         = ""
	staticDir        = ""
	frontendDir      = ""
	wallTemplate     *template.Template
	verifierTemplate *template.Template
	lotteryTemplate  *template.Template
)

func Init(acfg *config.AtomicConfig) {
	cfg := config.LoadCfgFromACfg(acfg)
	if cfg == nil {
		panic(errors.New("failed to get config file"))
	}
	imageDir = cfg.Common.ImageDir
	if _, err := os.Stat(imageDir); err != nil {
		panic(fmt.Errorf("dir %s not exist", imageDir))
	}
	staticDir = path.Join(cfg.Common.FrontendDir, "static")
	if _, err := os.Stat(staticDir); err != nil {
		panic(fmt.Errorf("dir %s not exist", staticDir))
	}
	frontendDir = cfg.Common.FrontendDir

	wallTemplate = template.Must(
		template.ParseFiles(path.Join(frontendDir, "wall.html")))
	verifierTemplate = template.Must(
		template.ParseFiles(path.Join(frontendDir, "verifier.html")))
	lotteryTemplate = template.Must(
		template.ParseFiles(path.Join(frontendDir, "lottery.html")))
}

func ServeWall(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	wallTemplate.Execute(w, r.Host)
}

func ServeVerifier(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	verifierTemplate.Execute(w, r.Host)
}

func ServeLottery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	lotteryTemplate.Execute(w, r.Host)
}
