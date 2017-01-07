package wechat

import (
	"WechatWall/backend/logger"
	//	"WechatWall/libredis"

	"github.com/chanxuehong/wechat.v2/mp/core"
	"github.com/chanxuehong/wechat.v2/mp/message/callback/request"
	"github.com/chanxuehong/wechat.v2/mp/message/callback/response"

	"net/http"
)

var log = logger.GetLogger("wechat")

const (
	wxAppId     = "wx9da3bbc39dab3cb8"
	wxAppSecret = "2184fa375f62ef0e46a17c05d2682d35"

	wxOriId         = "gh_d592b39f8508"
	wxToken         = "kevince"
	wxEncodedAESKey = "VgGi6gG0orzxL0J9x4qYXg95nNBblZgLWhdeXcbW3wK"
)

var (
	// 下面两个变量不一定非要作为全局变量, 根据自己的场景来选择.
	msgHandler core.Handler
	msgServer  *core.Server
)

func init() {
	mux := core.NewServeMux()
	mux.DefaultMsgHandleFunc(defaultMsgHandler)
	mux.MsgHandleFunc(request.MsgTypeText, textMsgHandler)

	mux.DefaultEventHandleFunc(defaultEventHandler)
	mux.EventHandleFunc(request.EventTypeSubscribe, subscribeEventHandler)
	mux.EventHandleFunc(request.EventTypeUnsubscribe, unsubscribeEventHandler)

	msgHandler = mux
	msgServer = core.NewServer(wxOriId, wxAppId, wxToken, wxEncodedAESKey, msgHandler, nil)
}

func textMsgHandler(ctx *core.Context) {
	log.Infof("received text message:\n%s\n", ctx.MsgPlaintext)

	msg := request.GetText(ctx.MixedMsg)
	resp := response.NewText(msg.FromUserName, msg.ToUserName, msg.CreateTime, msg.Content)
	//ctx.RawResponse(resp) // 明文回复
	ctx.AESResponse(resp, 0, "", nil) // aes密文回复
}

func defaultMsgHandler(ctx *core.Context) {
	log.Infof("received message:\n%s\n", ctx.MsgPlaintext)
	ctx.NoneResponse()
}

func subscribeEventHandler(ctx *core.Context) {

}

func unsubscribeEventHandler(ctx *core.Context) {

}

func defaultEventHandler(ctx *core.Context) {
	log.Infof("received event:\n%s\n", ctx.MsgPlaintext)
	ctx.NoneResponse()
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	msgServer.ServeHTTP(w, r, nil)
}
