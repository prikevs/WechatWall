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
	contentType := ""
	switch srcType {
	case "js":
		contentType = "text/javasrcipt; charset=utf8"
	case "css":
		contentType = "text/css; charset=utf-8"
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", contentType)
	log.Debug("path:", srcPath)
	http.ServeFile(w, r, srcPath)
}
