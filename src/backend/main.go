package main

import (
	"backend/wechat"

	"log"
	"net/http"
)

func init() {
	http.HandleFunc("/wx_callback", wechat.CallbackHandler)
}

func main() {
	log.Println(http.ListenAndServe(":9999", nil))
}
