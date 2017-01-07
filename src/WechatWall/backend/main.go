package main

import (
	"WechatWall/backend/wechat"

	"log"
	"net/http"
)

func init() {
	http.HandleFunc("/wx_callback", wechat.CallbackHandler)
}

func main() {
	log.Println(http.ListenAndServe("127.0.0.1:9999", nil))
}
