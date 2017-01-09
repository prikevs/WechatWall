package config

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"sync"
	"sync/atomic"
)

const (
	CONFIGJSON   = "backend-config.json"
	COMMONCONFIG = "common-config.json"
)

type CommonConfig struct {
	ImageDir    string `json:"image_dir"`
	FrontendDir string `json:"frontend_dir"`
	DebugF      bool
}

type LotteryConfig struct {
	// mode:
	// 0 all sent
	// 1 passed verification
	Mode int `json:"mode"`
}

type WechatConfig struct {
	WXAppId            string   `json:"wx_app_id"`
	WXAppSecret        string   `json:"wx_app_secret"`
	WXOriId            string   `json:"wx_ori_id"`
	WXToken            string   `json:"wx_token"`
	WXEncodedAESKey    string   `json:"wx_encoded_aes_key"`
	MessageOnSubscribe string   `json:"message_on_subscribe"`
	MessageOnReceived  string   `json:"message_on_received"`
	AdminList          []string `json:"admin_list"`
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
	Replay             bool `json:"replay"`
}

type Config struct {
	Common   CommonConfig
	Wechat   WechatConfig
	Verifier VerifierConfig
	Wall     WallConfig
	Lottery  LotteryConfig
}

type AtomicConfig struct {
	v  *atomic.Value
	mu sync.Mutex
}

func NewAtomicConfig(cfg *Config) *AtomicConfig {
	var v atomic.Value
	v.Store(cfg)
	return &AtomicConfig{
		v: &v,
	}
}

func (this *AtomicConfig) LoadConfig() *Config {
	return this.v.Load().(*Config)
}

func (this *AtomicConfig) StoreConfig(cfg Config) {
	this.mu.Lock()
	defer this.mu.Unlock()

	this.v.Store(&cfg)
}

func LoadCfgFromACfg(acfg *AtomicConfig) *Config {
	if acfg == nil {
		return nil
	}
	return acfg.LoadConfig()
}

func MustGetConfigJson(dir string) []byte {
	p := path.Join(dir, CONFIGJSON)
	data, err := ioutil.ReadFile(p)
	if err != nil {
		panic(err)
	}
	return data
}

func MustGetCommonConfigJson(dir string) []byte {
	p := path.Join(dir, COMMONCONFIG)
	data, err := ioutil.ReadFile(p)
	if err != nil {
		panic(err)
	}
	return data
}

func New(dir string) *Config {
	jconfig := MustGetConfigJson(dir)
	cconfig := MustGetCommonConfigJson(dir)
	c := &Config{}
	if err := json.Unmarshal([]byte(jconfig), c); err != nil {
		panic(err)
	}
	com := &CommonConfig{}
	if err := json.Unmarshal([]byte(cconfig), com); err != nil {
		panic(err)
	}
	c.Common = *com
	return c
}

func GetConfig(acfg *AtomicConfig, dft interface{},
	getter func(cfg *Config) interface{}) interface{} {

	cfg := LoadCfgFromACfg(acfg)
	if cfg == nil {
		return dft
	}
	return getter(cfg)
}
