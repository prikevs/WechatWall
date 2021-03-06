package wechat

import (
	"WechatWall/backend/config"
	"WechatWall/libredis"
	"WechatWall/logger"

	"github.com/chanxuehong/wechat.v2/mp/core"
	"github.com/chanxuehong/wechat.v2/mp/message/callback/request"
	"github.com/chanxuehong/wechat.v2/mp/message/callback/response"

	"net/http"
	"time"
)

var log = logger.GetLogger("backend/wechat")

var (
	wxAppId     = ""
	wxAppSecret = ""

	wxOriId         = ""
	wxToken         = ""
	wxEncodedAESKey = ""

	MessageOnSubscribe = ""
	MessageOnReceived  = ""
)

var (
	msgHandler core.Handler
	msgServer  *core.Server

	pMQ     libredis.MQ
	vMQ     libredis.MQ
	okSet   libredis.Set
	wxSet   libredis.Set
	sentSet libredis.Set
)

func Init(acfg *config.AtomicConfig) {
	InitConsole(acfg)

	cfg := config.LoadCfgFromACfg(acfg)
	if cfg != nil {
		wxAppId = cfg.Wechat.WXAppId
		wxAppSecret = cfg.Wechat.WXAppSecret
		wxOriId = cfg.Wechat.WXOriId
		wxToken = cfg.Wechat.WXToken
		wxEncodedAESKey = cfg.Wechat.WXEncodedAESKey
		MessageOnSubscribe = cfg.Wechat.MessageOnSubscribe
		MessageOnReceived = cfg.Wechat.MessageOnReceived
	}

	mux := core.NewServeMux()
	mux.DefaultMsgHandleFunc(defaultMsgHandler)
	mux.MsgHandleFunc(request.MsgTypeText, textMsgHandler)

	mux.DefaultEventHandleFunc(defaultEventHandler)
	mux.EventHandleFunc(request.EventTypeSubscribe, subscribeEventHandler)
	mux.EventHandleFunc(request.EventTypeUnsubscribe, unsubscribeEventHandler)
	msgHandler = mux
	msgServer = core.NewServer(wxOriId, wxAppId, wxToken, wxEncodedAESKey, msgHandler, nil)
}

func FailOnError(err error) {
	if err != nil {
		log.Critical("Failed to init libredis", err)
		panic(err)
	}
}

func init() {
	var err error

	pMQ, err = libredis.GetPMQ()
	FailOnError(err)
	vMQ, err = libredis.GetVMQ()
	FailOnError(err)

	okSet, err = libredis.GetOKSet()
	FailOnError(err)
	wxSet, err = libredis.GetWXSet()
	FailOnError(err)
	sentSet, err = libredis.GetSentSet()
	FailOnError(err)
}

func textMsgHandler(ctx *core.Context) {
	msg := request.GetText(ctx.MixedMsg)
	if len(msg.Content) > 0 && msg.Content[0] == ':' {
		log.Infof("received commands from %s: %s", msg.FromUserName, msg.Content)
		result := handleCommand(msg.FromUserName, msg.Content[1:])
		resp := response.NewText(msg.FromUserName, msg.ToUserName, msg.CreateTime, result)
		ctx.AESResponse(resp, 0, "", nil)
		return
	}

	log.Infof("received text message from %s: %s", msg.FromUserName, msg.Content)
	rmsg := &libredis.Msg{
		UserOpenid: msg.FromUserName,
		CreateTime: msg.CreateTime,
		Content:    msg.Content,
		MsgId:      ctx.MixedMsg.MsgId,
		MsgType:    "text",
		TTL:        libredis.MsgTTL,
		AddTime:    time.Now(),
	}

	var err error
	// Add user to sent set
	_, err = sentSet.Add(msg.FromUserName)
	FailOnError(err)

	// Add msg to pmq
	err = libredis.PublishClassToMQ(rmsg, pMQ)
	FailOnError(err)

	// TODO: Add suitable response
	resp := response.NewText(msg.FromUserName, msg.ToUserName, msg.CreateTime, MessageOnReceived)
	ctx.RawResponse(resp) // 明文回复
	//ctx.AESResponse(resp, 0, "", nil) // aes密文回复
}

func defaultMsgHandler(ctx *core.Context) {
	log.Infof("received message:\n%s\n", ctx.MsgPlaintext)
	ctx.NoneResponse()
}

func subscribeEventHandler(ctx *core.Context) {
	var err error
	event := request.GetSubscribeEvent(ctx.MixedMsg)
	log.Info("received subscribe event", event.MsgHeader)

	header := event.MsgHeader
	ismem, err := okSet.IsMember(header.FromUserName)
	FailOnError(err)
	if !ismem {
		_, err = wxSet.Add(header.FromUserName)
		FailOnError(err)
	}
	// add welcome message
	ctx.NoneResponse()
}

func unsubscribeEventHandler(ctx *core.Context) {

}

func defaultEventHandler(ctx *core.Context) {
	log.Infof("received event:\n%s\n", ctx.MsgPlaintext)
	ctx.NoneResponse()
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	// log.Debug(*r)
	msgServer.ServeHTTP(w, r, nil)
}
