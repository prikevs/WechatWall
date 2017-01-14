package lottery

import (
	"WechatWall/backend/config"
	"WechatWall/backend/utils"
	"WechatWall/libredis"
	"WechatWall/logger"

	"github.com/gorilla/mux"

	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

var log = logger.GetLogger("backend/lottery")

var (
	ACfg  *config.AtomicConfig
	dMode = 1

	okSet    libredis.Set
	passSet  libredis.Set
	sentSet  libredis.Set
	usersMap libredis.Map
	lvmMap   libredis.Map
)

func LoadMode() int {
	return config.GetConfig(ACfg, dMode,
		func(cfg *config.Config) interface{} {
			switch cfg.Lottery.Mode {
			case 0:
				return 0
			case 1:
				return 1
			case 2:
				return 2
			default:
				return dMode
			}
		}).(int)
}

func LoadBackDoor() string {
	return config.GetConfig(ACfg, "",
		func(cfg *config.Config) interface{} {
			return cfg.Lottery.BackDoor
		}).(string)
}

func ResetBackDoor() bool {
	return config.SetConfig(ACfg, "",
		func(cfg *config.Config, val interface{}) {
			cfg.Lottery.BackDoor = val.(string)
		})
}

func FailOnError(err error) {
	if err != nil {
		log.Critical(err)
		panic(err)
	}
}

func init() {
	var err error
	okSet, err = libredis.GetOKSet()
	FailOnError(err)
	passSet, err = libredis.GetPassSet()
	FailOnError(err)
	sentSet, err = libredis.GetSentSet()
	FailOnError(err)
	usersMap, err = libredis.GetUsersMap()
	FailOnError(err)
	lvmMap, err = libredis.GetLVMMap()
	FailOnError(err)
}

func Init(acfg *config.AtomicConfig) {
	ACfg = acfg
}

type UserInfo struct {
	Username   string `json:"username"`
	UserOpenid string `json:"user_openid"`
	Msg        string `json:"msg"`
	ImgUrl     string `json:"img_url"`
}

type RespMsg struct {
	RetCode  int        `json:"ret_code"`
	ErrMsg   string     `json:"err_msg"`
	UserList []UserInfo `json:"user_list"`
}

type BDMsg struct {
	HasBD  bool   `json:"has_bd"`
	Openid string `json:"openid"`
}

func GetLotteryOpenids(mode int) ([]string, error) {
	var targetSet libredis.Set
	switch mode {
	case 0:
		targetSet = sentSet
	case 1:
		targetSet = passSet
	case 2:
		targetSet = okSet
	}
	resSlice, err := okSet.Inter(targetSet)
	if err != nil {
		return nil, err
	}
	for i := range resSlice {
		j := rand.Intn(i + 1)
		resSlice[i], resSlice[j] = resSlice[j], resSlice[i]
	}
	return resSlice, nil
}

func BuildUserInfoByOpenid(mode int, id string) (*UserInfo, error) {
	user := &libredis.User{}
	if err := libredis.GetClassFromMap(id, user, usersMap); err != nil {
		return nil, err
	}
	// get message from map
	msg := ""
	switch mode {
	case 1:
		var err error
		msg, err = lvmMap.Get(id)
		if err != nil {
			return nil, err
		}
	case 0:
		msg = ""
	}
	return &UserInfo{
		Username:   user.UserName,
		UserOpenid: user.UserOpenid,
		Msg:        msg,
		ImgUrl:     utils.BuildImagePath(id),
	}, nil
}

func GetUserInfos() ([]UserInfo, error) {
	mode := LoadMode()
	ids, err := GetLotteryOpenids(mode)
	if err != nil {
		log.Error("failed to get lottery user openids:", err)
	}
	users := make([]UserInfo, 0)
	for _, id := range ids {
		user, err := BuildUserInfoByOpenid(mode, id)
		if err != nil {
			log.Warning("failed to get build message for", id, "errmsg:", err)
		}
		users = append(users, *user)
	}
	return users, nil
}

func WriteResponse(w http.ResponseWriter, r *http.Request, resp interface{}) {
	data, err := json.Marshal(resp)
	if err != nil {
		log.Error("failed to encode data to json:", err)
		log.Debug("data:", resp)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf8")
	callback := r.URL.Query().Get("callback")

	if callback != "" {
		fmt.Fprintf(w, "%s(%s)", callback, data)
	} else {
		w.Write(data)
	}
}

func GenResponse(w http.ResponseWriter, r *http.Request) {
	resp := &RespMsg{RetCode: 200}
	defer WriteResponse(w, r, resp)

	users, err := GetUserInfos()
	if err != nil {
		log.Error("failed to generate response:", err)
		resp.RetCode = 500
		resp.ErrMsg = err.Error()
		return
	}
	resp.UserList = users
}

func ServeLotteryFunc(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	funcName := vars["func"]
	// log.Debug(*r)
	switch funcName {
	case "list":
		log.Warning("here comes a request to lottery api, PAY ATTENTION!")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		GenResponse(w, r)
	case "bd":
		log.Warning("here comes a request to BACK DOOR API")
		bdmsg := &BDMsg{}
		var bd string
		bd = LoadBackDoor()
		if bd == "" {
			bdmsg.HasBD = false
		} else {
			bdmsg.HasBD = true
			bdmsg.Openid = bd
			if ResetBackDoor() == false {
				log.Warning("failed to reset back door")
			}
		}
		WriteResponse(w, r, bdmsg)
	}
}
