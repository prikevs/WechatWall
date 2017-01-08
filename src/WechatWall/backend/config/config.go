package config

import (
	"encoding/json"
	"io/ioutil"
	"path"
)

const (
	CONFIGJSON = "backend-config.json"
)

type WechatConfig struct {
	WXAppId         string `json:"wx_app_id"`
	WXAppSecret     string `json:"wx_app_secret"`
	WXOriId         string `json:"wx_ori_id"`
	WXToken         string `json:"wx_token"`
	WXEncodedAESKey string `json:"wx_encoded_aes_key"`
}

type VerifierConfig struct {
	ReadyPipSize             int    `json:"ready_pip_size"`
	MaxMsgWaitingTime        int    `json:"max_msg_waiting_time"`
	SendVerificationDuration int    `json:"send_verification_duration"`
	SendNotification         bool   `json:"send_notification"`
	NotificationMessage      string `json:"notification_message"`
	StrictOrigin             bool   `json:"strict_origin"`
	NeedVerification         bool   `json:"need_verification"`
}

type WallConfig struct {
	SendToWallDuration int  `json:"send_to_wall_duration"`
	ReliableMsg        bool `json:"reliable_msg"`
}

type Config struct {
	Wechat   WechatConfig
	Verifier VerifierConfig
	Wall     WallConfig
}

func MustGetConfigJson(dir string) []byte {
	p := path.Join(dir, CONFIGJSON)
	data, err := ioutil.ReadFile(p)
	if err != nil {
		panic(err)
	}
	return data
}

func New(dir string) *Config {
	jconfig := MustGetConfigJson(dir)
	c := &Config{}
	if err := json.Unmarshal([]byte(jconfig), c); err != nil {
		panic(err)
	}
	return c
}
