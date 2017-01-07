package main

import (
	"WechatWall/backend/verifier"
	"WechatWall/backend/wall"
	"WechatWall/backend/wechat"

	"log"
	"net/http"
)

func init() {
	http.HandleFunc("/wx_callback", wechat.CallbackHandler)
	http.HandleFunc("/ws/verifier", verifier.ServeVerifierWS)
	http.HandleFunc("/ws/wall", wall.ServeWallWS)
}

func main() {
	log.Println(http.ListenAndServe("127.0.0.1:9999", nil))
}
