package web

import (
	"github.com/gorilla/mux"

	"net/http"
	"path"
)

func ServeStatic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	srcType := vars["type"]
	srcName := vars["name"]
	srcPath := path.Join(staticDir, srcType, srcName)
	log.Debug("request to", srcType, "/", srcName)
	log.Debug("path:", srcPath)
	http.ServeFile(w, r, srcPath)
}
