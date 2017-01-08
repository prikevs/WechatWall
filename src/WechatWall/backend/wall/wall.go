package wall

/*
{
	"msg_id": xxxxxx,
	"username": "Kevince",
	"openid": "xxxxxx",
	"msg_type": "text",
	"create_time": 12324353,
	"content": "恭喜恭喜！",
	"img_url": "/img/xxxxx.jpg"
}
*/

import (
	"WechatWall/backend/config"
	"WechatWall/libredis"
	"WechatWall/logger"

	"net/http"
	"time"
)

var log = logger.GetLogger("backend/wall")

// config
const (
	TickWarnRound = 3
)

var (
	SendToWallDuration = 2 * time.Second
	ReliableMsg        = true
)

var (
	hub *Hub
	vMQ libredis.MQ
)

func FailOnError(err error) {
	if err != nil {
		log.Critical(err)
		panic(err)
	}
}

func init() {
	var err error
	vMQ, err = libredis.GetVMQ()
	FailOnError(err)

}

func Init(cfg *config.WallConfig) {
	SendToWallDuration = time.Duration(cfg.SendToWallDuration) * time.Second
	ReliableMsg = cfg.ReliableMsg

	bc := make(chan bool)
	wallmsgs := make(chan libredis.Msg)
	hub = newHub(wallmsgs, bc)
	go hub.run()
	go tickWallBroadcastSignal(bc, SendToWallDuration)
	go prepareWallMsgs(wallmsgs)
}

func prepareWallMsgs(wallmsgs chan<- libredis.Msg) {
	for {
		msg := &libredis.Msg{}
		if err := libredis.ConsumeClassFromMQ(msg, vMQ, 0); err != nil {
			log.Error("failed to get message from verified mq:", err)
			continue
		}
		wallmsgs <- *msg
	}
}

type WallMsg struct {
	MsgId      int64 `json:"msg_id"`
	Username   string
	Openid     string
	MsgType    string `json:"msg_type"`
	CreateTime int64  `json:"create_time"`
	Content    string
	ImgUrl     string `json:"img_url"`
}

func tickWallBroadcastSignal(bc chan bool, d time.Duration) {
	round := 0
	for range time.Tick(d) {
		select {
		case bc <- true:
			round = 0
		default:
			round++
			if round >= TickWarnRound {
				log.Warning("wall tick block",
					round,
					"round(s), maybe hub is doing some heavy work, reset tick")
				round = 0
			}
		}
	}
}

func ServeWallWS(w http.ResponseWriter, r *http.Request) {
	serveWs(hub, w, r)
}
