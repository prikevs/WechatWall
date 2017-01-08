package imager

import (
	"WechatWall/backend/config"
	"WechatWall/logger"

	"github.com/gorilla/mux"

	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
)

var log = logger.GetLogger("backend/imager")

var (
	imageDir = ""
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
}

func ServeImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageName := vars["image"]
	log.Debug("request to", imageName)
	imagePath := path.Join(imageDir, imageName)
	w.Header().Set("Content-Type", "image/jpg")
	http.ServeFile(w, r, imagePath)
}
