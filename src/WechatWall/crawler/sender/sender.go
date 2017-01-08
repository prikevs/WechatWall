package sender

import (
	"WechatWall/crawler/config"
	"WechatWall/libredis"
	"WechatWall/logger"

	"encoding/json"
)

var log = logger.GetLogger("crawler/sender")

var (
	smq libredis.MQ
)

type Sender struct {
	cfg    *config.Config
	poster *Poster
}

func FailOnError(err error) {
	if err != nil {
		log.Critical(err)
		panic(err)
	}
}

func init() {
	var err error
	smq, err = libredis.GetSMQ()
	FailOnError(err)
}

func NewSender(cfg *config.Config) *Sender {
	return &Sender{
		cfg:    cfg,
		poster: NewPoster(cfg),
	}
}

type Resp struct {
	BaseResp struct {
		ErrMsg string `json:"err_msg"`
		Ret    int
	} `json:"base_resp"`
}

func (this *Sender) Send(msg libredis.Msg) {
	data, err := this.poster.Do(&msg)
	if err != nil {
		log.Error("failed to send msg to", msg.UserOpenid, ":", err)
		return
	}
	resp := &Resp{}
	if err := json.Unmarshal(data, resp); err != nil {
		log.Error("failed to parse response:", err, "rawmsg:", string(data))
		return
	}
	if resp.BaseResp.ErrMsg != "" || resp.BaseResp.Ret != 0 {
		log.Error("failed to send msg to", msg.UserOpenid,
			"due to return code err: err_msg:", string(data))
		return
	}
	log.Info("message sent to", msg.UserOpenid, msg.Username)
}

func Run(cfg *config.Config) {
	log.Info("starting to run sender loop")
	sender := NewSender(cfg)
	for {
		msg := &libredis.Msg{}
		if err := libredis.ConsumeClassFromMQ(msg, smq, 0); err != nil {
			log.Error("cannot consume msg from smq,", err)
			continue
		}
		go sender.Send(*msg)
	}
}
