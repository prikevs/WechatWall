package web

import (
	"github.com/gorilla/mux"

	"net/http"
	"path"
)

func ServeImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageName := vars["image"]
	log.Debug("request to", imageName)
	imagePath := path.Join(imageDir, imageName)
	w.Header().Set("Content-Type", "image/jpg")
	http.ServeFile(w, r, imagePath)
}
