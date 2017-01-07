package verifier

/*
{
	"username": "Kevince",
	"openid": "o_5d9s-is4Y2UWt7u-CC4UMDYsbE",
	"msg_id": 6372734764678826899,
	"msg_type": "text",
	"create_time": 1483767937,
	"content": "恭喜",
	"img_url": "/img/xxxxx.jpg"
}

{
	"msg_id": xxxxxx,
	"verified_time": 1483767999,
	"show_now": true/false,
}
*/

import (
	"WechatWall/libredis"
	"WechatWall/logger"

	"net/http"
	"time"
)

var log = logger.GetLogger("verifier")

const (
	ReadyPipeSize     = 20
	MaxMsgWaitingTime = 5 * 60 * time.Second
)

var (
	hub *Hub

	pMQ      libredis.MQ
	vMQ      libredis.MQ
	okSet    libredis.Set
	pMsgsMap libredis.Map
)

func FailOnError(err error) {
	if err != nil {
		log.Critical(err)
		panic(err)
	}
}

func init() {
	// init libredis
	var err error
	pMQ, err = libredis.GetPMQ()
	FailOnError(err)
	vMQ, err = libredis.GetVMQ()
	FailOnError(err)
	okSet, err = libredis.GetOKSet()
	FailOnError(err)
	pMsgsMap, err = libredis.GetPMsgsMap()
	FailOnError(err)

	readymsgs := make(chan libredis.Msg, ReadyPipeSize)
	hub = newHub(readymsgs)
	go hub.run()
	go prepareMsgs(readymsgs)
}

func prepareMsgs(readymsgs chan<- libredis.Msg) {
	for {
		msg := &libredis.Msg{}
		if err := libredis.ConsumeClassFromMQ(msg, pMQ, 0); err != nil {
			log.Error("failed to get message from pending mq:", err)
			continue
		}
		ismem, err := okSet.IsMember(msg.UserOpenid)
		if err != nil {
			log.Error("failed to check if user in ok set:", err)
			continue
		}
		if !ismem {
			log.Debugf("msg from %s will be republish due to user not in ok set",
				msg.UserOpenid)
			if time.Since(msg.AddTime).Seconds() > msg.TTL.Seconds() {
				log.Warningf("msg from %s is discarded due to TTL", msg.UserOpenid)
				continue
			}
			if err := libredis.PublishClassToMQ(msg, pMQ); err != nil {
				log.Errorf("failed to republis msg %s to pending mq: %s",
					msg.UserOpenid, err.Error())
			}
		} else {
			log.Debugf("msg from %s is ready to be verified", msg.UserOpenid)
			readymsgs <- *msg
		}
	}
}

func serveVerifier(w http.ResponseWriter, r *http.Request) {
	serveWs(hub, w, r)
}
