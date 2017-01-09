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
	"WechatWall/backend/config"
	"WechatWall/libredis"
	"WechatWall/logger"

	"net/http"
	"time"
)

var log = logger.GetLogger("backend/verifier")

// some config variables
var (
	ReadyPipSize        = 20
	TickWarnRound       = 3
	NotificationMessage = "消息已通过审核"
	StrictOrigin        = false

	dMaxMsgWaitingTime        = 5 * 60 * time.Second
	dSendVerificationDuration = 1 * time.Second
	dSendNotification         = true
	dNeedVerification         = true

	ACfg *config.AtomicConfig
)

func LoadNeedVerification() bool {
	cfg := config.LoadCfgFromACfg(ACfg)
	if cfg == nil {
		return dNeedVerification
	}
	return cfg.Verifier.NeedVerification
}

func LoadSendNotification() bool {
	cfg := config.LoadCfgFromACfg(ACfg)
	if cfg == nil {
		return dSendNotification
	}
	return cfg.Verifier.SendNotification
}

func LoadSendVerificationDuration() time.Duration {
	cfg := config.LoadCfgFromACfg(ACfg)
	if cfg == nil {
		return dSendVerificationDuration
	}
	if cfg.Verifier.SendVerificationDuration == 0 {
		return dSendVerificationDuration
	}
	return time.Duration(cfg.Verifier.SendVerificationDuration) * time.Second
}

func LoadMaxMsgWaitingTime() time.Duration {
	cfg := config.LoadCfgFromACfg(ACfg)
	if cfg == nil {
		return dMaxMsgWaitingTime
	}
	if cfg.Verifier.MaxMsgWaitingTime == 0 {
		return dMaxMsgWaitingTime
	}
	return time.Duration(cfg.Verifier.MaxMsgWaitingTime) * time.Second
}

var (
	hub *Hub

	pMQ      libredis.MQ
	vMQ      libredis.MQ
	sMQ      libredis.MQ
	okSet    libredis.Set
	passSet  libredis.Set
	lvmMap   libredis.Map
	pMsgsMap libredis.Map
	usersMap libredis.Map
)

// verification message to send
type VMsgSent struct {
	Username   string `json:"username"`
	Openid     string `json:"openid"`
	MsgId      string `json:"msg_id"`
	MsgType    string `json:"msg_type"`
	CreateTime int64  `json:"create_time"`
	Content    string `json:"content"`
	ImgUrl     string `json:"img_url"`
	TTL        int64  `json:"ttl"`
}

// verification message to receive
type VMsgRecvd struct {
	MsgId        string `json:"msg_id"`
	VerifiedTime int64  `json:"verified_time"`
	ShowNow      bool   `json:"show_now"`
}

type VMsgResp struct {
	MsgId   string `json:"msg_id"`
	RetCode int    `json:"ret_code"`
	ErrMsg  string `json:"err_msg"`
}

func FailOnError(err error) {
	if err != nil {
		log.Critical(err)
		panic(err)
	}
}

func tickBroadcastSignal(bc chan bool, d time.Duration) {
	round := 0
	for range time.Tick(d) {
		select {
		case bc <- true:
			round = 0
		default:
			round++
			if round >= TickWarnRound {
				log.Warning("wall tick not working for",
					round,
					"round(s), maybe hub is doing some heavy work, reset tick")
				round = 0
			}
		}
	}
}

func init() {
	// init libredis
	var err error
	pMQ, err = libredis.GetPMQ()
	FailOnError(err)
	vMQ, err = libredis.GetVMQ()
	FailOnError(err)
	sMQ, err = libredis.GetSMQ()
	FailOnError(err)
	okSet, err = libredis.GetOKSet()
	FailOnError(err)
	passSet, err = libredis.GetPassSet()
	FailOnError(err)
	pMsgsMap, err = libredis.GetPMsgsMap()
	FailOnError(err)
	usersMap, err = libredis.GetUsersMap()
	FailOnError(err)
	lvmMap, err = libredis.GetLVMMap()
	FailOnError(err)
}

func Init(acfg *config.AtomicConfig) {
	ACfg = acfg

	cfg := config.LoadCfgFromACfg(acfg)
	if cfg != nil {
		ReadyPipSize = cfg.Verifier.ReadyPipSize
		NotificationMessage = cfg.Verifier.NotificationMessage
		StrictOrigin = cfg.Verifier.StrictOrigin
	}

	readymsgs := make(chan libredis.Msg, ReadyPipSize)
	bc := make(chan bool)
	hub = newHub(readymsgs, bc)
	go hub.run()
	go prepareMsgs(readymsgs)
	go tickBroadcastSignal(bc, LoadSendVerificationDuration())
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
			// get user info
			user := &libredis.User{}
			if err := libredis.GetClassFromMap(msg.UserOpenid, user, usersMap); err != nil {
				log.Error("failed to get user info from users map")
				continue
			}
			msg.Username = user.UserName

			log.Debugf("msg from %s is ready to be verified", msg.Username)
			readymsgs <- *msg
		}
	}
}

func ServeVerifierWS(w http.ResponseWriter, r *http.Request) {
	serveWs(hub, w, r)
}
