package wechat

import (
	"WechatWall/backend/config"
	"WechatWall/libredis"

	"errors"
	"fmt"
	"strconv"
	"strings"
)

type UserConsole struct {
	vMQ      libredis.MQ
	pMQ      libredis.MQ
	owSet    libredis.Set
	passSet  libredis.Set
	pMsgsMap libredis.Map
}

type AdminConsole struct {
	UserConsole
	adminList []string
	acfg      *config.AtomicConfig
}

func (this *AdminConsole) IsAdmin(openid string) bool {
	for _, v := range this.adminList {
		if openid == v {
			return true
		}
	}
	return false
}

func (this *UserConsole) getVMQLength() (length int64, err error) {
	length, err = this.vMQ.Length()
	return
}

func (this *UserConsole) getPMQLength() (length int64, err error) {
	length, err = this.pMQ.Length()
	return
}

func (this *UserConsole) getNumberOfWallMsgs() (num int64, err error) {
	return this.owSet.Size()
}
func (this *UserConsole) getNumberOfPassedUsers() (num int64, err error) {
	return this.passSet.Size()
}
func (this *UserConsole) getNumberOfWaitingMsgs() (num int64, err error) {
	return this.pMsgsMap.Size()
}

var (
	adminCon *AdminConsole
	userCon  *UserConsole
)

func InitConsole(acfg *config.AtomicConfig) {
	cfg := config.LoadCfgFromACfg(acfg)
	if cfg == nil {
		FailOnError(errors.New("config not loaded"))
	}

	vMQ, err := libredis.GetVMQ()
	FailOnError(err)
	pMQ, err := libredis.GetPMQ()
	FailOnError(err)
	owSet, err := libredis.GetOWSet()
	FailOnError(err)
	passSet, err := libredis.GetPassSet()
	FailOnError(err)
	pMsgsMap, err := libredis.GetPMsgsMap()
	FailOnError(err)

	userCon := &UserConsole{
		vMQ:      vMQ,
		pMQ:      pMQ,
		owSet:    owSet,
		passSet:  passSet,
		pMsgsMap: pMsgsMap,
	}

	var adminList = make([]string, len(cfg.Wechat.AdminList))
	copy(adminList, cfg.Wechat.AdminList)
	adminCon = &AdminConsole{
		UserConsole: *userCon,
		adminList:   adminList,
		acfg:        acfg,
	}
}

type Console interface {
	CmdQuery([]string) string
	CmdSet([]string) string
	CmdHelp([]string) string
}

const (
	ErrInvalidCmd    = "invalid command"
	ErrUnauthorized  = "unauthorized usage"
	ErrInternalError = "internal error"
	Success          = "command execute successfully"
)

func SetConfig(acfg *config.AtomicConfig, val interface{}, setter func(cfg *config.Config, val interface{})) string {
	cfg := config.LoadCfgFromACfg(acfg)
	if cfg == nil {
		return ErrInternalError
	}
	setter(cfg, val)
	acfg.StoreConfig(*cfg)
	return Success
}

func (this *AdminConsole) SetLMode(cmds []string) string {
	if len(cmds) != 1 {
		return ErrInvalidCmd
	}
	lmode, err := strconv.Atoi(cmds[0])
	if lmode != 0 && lmode != 1 {
		return ErrInvalidCmd
	}
	if err != nil {
		return ErrInvalidCmd
	}
	return SetConfig(this.acfg, lmode,
		func(cfg *config.Config, val interface{}) {
			cfg.Lottery.Mode = val.(int)
		})
}

func (this *AdminConsole) SetSendn(cmds []string) string {
	if len(cmds) != 1 {
		return ErrInvalidCmd
	}
	switch cmds[0] {
	case "on":
		log.Info("ADMIN set send notification on")
		return SetConfig(this.acfg, true,
			func(cfg *config.Config, val interface{}) {
				cfg.Verifier.SendNotification = val.(bool)
			})
	case "off":
		log.Info("ADMIN set send notification off")
		return SetConfig(this.acfg, false,
			func(cfg *config.Config, val interface{}) {
				cfg.Verifier.SendNotification = val.(bool)
			})
	default:
		return ErrInvalidCmd
	}
}

func (this *AdminConsole) SetNeedv(cmds []string) string {
	if len(cmds) != 1 {
		return ErrInvalidCmd
	}
	switch cmds[0] {
	case "on":
		log.Info("ADMIN set need verification on")
		return SetConfig(this.acfg, true,
			func(cfg *config.Config, val interface{}) {
				cfg.Verifier.NeedVerification = val.(bool)
			})
	case "off":
		log.Info("ADMIN set send verification off")
		return SetConfig(this.acfg, false,
			func(cfg *config.Config, val interface{}) {
				cfg.Verifier.NeedVerification = val.(bool)
			})
	default:
		return ErrInvalidCmd
	}
}

func (this *AdminConsole) SetTTL(cmds []string) string {
	if len(cmds) != 1 {
		return ErrInvalidCmd
	}
	ttl, err := strconv.Atoi(cmds[0])
	if err != nil {
		return ErrInvalidCmd
	}
	if ttl <= 0 {
		return ErrInvalidCmd
	}
	return SetConfig(this.acfg, ttl,
		func(cfg *config.Config, val interface{}) {
			cfg.Verifier.MaxMsgWaitingTime = val.(int)
		})
}

func (this *AdminConsole) SetReplay(cmds []string) string {
	if len(cmds) != 1 {
		return ErrInvalidCmd
	}
	switch cmds[0] {
	case "on":
		log.Info("ADMIN set need verification on")
		return SetConfig(this.acfg, true,
			func(cfg *config.Config, val interface{}) {
				cfg.Wall.Replay = val.(bool)
			})
	case "off":
		log.Info("ADMIN set send verification off")
		return SetConfig(this.acfg, false,
			func(cfg *config.Config, val interface{}) {
				cfg.Wall.Replay = val.(bool)
			})
	default:
		return ErrInvalidCmd
	}
}

func (this *AdminConsole) SetSwd(cmds []string) string {
	if len(cmds) != 1 {
		return ErrInvalidCmd
	}
	swd, err := strconv.Atoi(cmds[0])
	if swd <= 0 {
		return ErrInvalidCmd
	}
	if err != nil {
		return ErrInvalidCmd
	}
	return SetConfig(this.acfg, swd,
		func(cfg *config.Config, val interface{}) {
			cfg.Wall.SendToWallDuration = val.(int)
		})
}

func (this *AdminConsole) SetSvd(cmds []string) string {
	if len(cmds) != 1 {
		return ErrInvalidCmd
	}
	svd, err := strconv.Atoi(cmds[0])
	if svd <= 0 {
		return ErrInvalidCmd
	}
	if err != nil {
		return ErrInvalidCmd
	}
	return SetConfig(this.acfg, svd,
		func(cfg *config.Config, val interface{}) {
			cfg.Verifier.SendVerificationDuration = val.(int)
		})
}

func (this *AdminConsole) CmdSet(cmds []string) string {
	if len(cmds) == 0 {
		return ErrInvalidCmd
	}
	switch cmds[0] {
	case "lmode":
		return this.SetLMode(cmds[1:])
	case "sendn":
		return this.SetSendn(cmds[1:])
	case "needv":
		return this.SetNeedv(cmds[1:])
	case "ttl":
		return this.SetTTL(cmds[1:])
	case "replay":
		return this.SetReplay(cmds[1:])
	case "swd":
		return this.SetSwd(cmds[1:])
	case "svd":
		return this.SetSvd(cmds[1:])
	}
	return ErrInvalidCmd
}

func (this *UserConsole) CmdSet(cmds []string) string {
	return ErrUnauthorized
}

func (this *AdminConsole) CmdHelp(cmds []string) string {
	return this.UserConsole.CmdHelp(cmds) + `Admin Console Help:
query/q/?
lmode:
   lottery mode, 1: all sent users, 2: passed users
ttl:
    current ttl of unverified message
svd:
    current duration of sending message to verification
swd:
    current duration of sending message to wall
sendn:
    current status of if sending notification
needv:
    current status of if needing verification

set/s
lmode:
   lottery mode, 1: all sent users, 2: passed users
ttl <int>:
    set ttl of unverified message (in seconds)
svd <int>:
    set duration of sending message to verification (in seconds)
swd <int>:
    set duration of sending message to wall (in seconds)
sendn <on/off>:
    set sending notification on/off
needv <on/off>:
    set needing verification on/off
`
}

func (this *UserConsole) CmdHelp(cmds []string) string {
	return `Console Help:
query/q/?
vmq:
    length of MQ for verified messages
pmq:
    length of MQ for pending messages
wmsg:
    number of messages sent to verifier but not verified
ow:
    number of messages already on wall
pass:
    number of users have messages passed verification

`
}

func GetConfig(acfg *config.AtomicConfig, getter func(cfg *config.Config) interface{}) interface{} {
	cfg := config.LoadCfgFromACfg(acfg)
	if cfg == nil {
		return ErrInternalError
	}
	return getter(cfg)
}

func (this *AdminConsole) CmdQuery(cmds []string) string {
	if len(cmds) == 0 {
		return ErrInvalidCmd
	}
	result := this.UserConsole.CmdQuery(cmds)
	if result == ErrInvalidCmd {
		switch cmds[0] {
		case "lmode":
			mode := GetConfig(this.acfg,
				func(cfg *config.Config) interface{} {
					return cfg.Lottery.Mode
				}).(int)
			result = fmt.Sprintf("lottery mode is: %d", mode)

		case "ttl":
			ttl := GetConfig(this.acfg,
				func(cfg *config.Config) interface{} {
					return cfg.Verifier.MaxMsgWaitingTime
				}).(int)
			result = fmt.Sprintf("current TTL is: %d", ttl)
		case "swd":
			swd := GetConfig(this.acfg,
				func(cfg *config.Config) interface{} {
					return cfg.Wall.SendToWallDuration
				}).(int)
			result = fmt.Sprintf("current send to wall duration is: %d", swd)

		case "svd":
			svd := GetConfig(this.acfg,
				func(cfg *config.Config) interface{} {
					return cfg.Verifier.SendVerificationDuration
				}).(int)
			result = fmt.Sprintf("current send to verifer duration is: %d", svd)

		case "sendn":
			sendn := GetConfig(this.acfg,
				func(cfg *config.Config) interface{} {
					return cfg.Verifier.SendNotification
				}).(bool)
			mode := ""
			if sendn {
				mode = "on"
			} else {
				mode = "off"
			}
			result = fmt.Sprintf("sendn mode %v", mode)
		case "needv":
			needv := GetConfig(this.acfg,
				func(cfg *config.Config) interface{} {
					return cfg.Verifier.NeedVerification
				}).(bool)
			mode := ""
			if needv {
				mode = "on"
			} else {
				mode = "off"
			}
			result = fmt.Sprintf("needv mode %v", mode)
		case "df":
			df := GetConfig(this.acfg,
				func(cfg *config.Config) interface{} {
					return cfg.Common.DebugF
				}).(bool)
			result = fmt.Sprintf("debug mode %v", df)
		case "replay":
			replay := GetConfig(this.acfg,
				func(cfg *config.Config) interface{} {
					return cfg.Wall.Replay
				}).(bool)
			mode := ""
			if replay {
				mode = "on"
			} else {
				mode = "off"
			}
			result = fmt.Sprintf("replay mode %v", mode)
		}
	}
	return result
}

func (this *UserConsole) CmdQuery(cmds []string) string {
	if len(cmds) == 0 {
		return ErrInvalidCmd
	}
	switch cmds[0] {
	case "vmq":
		length, err := this.getVMQLength()
		if err != nil {
			return err.Error()
		}
		return fmt.Sprintf("len mq:vmq: %d", length)
	case "pmq":
		length, err := this.getPMQLength()
		if err != nil {
			return err.Error()
		}
		return fmt.Sprintf("len mq:pmq: %d", length)
	case "wmsg":
		total, err := this.getNumberOfWaitingMsgs()
		if err != nil {
			return err.Error()
		}
		return fmt.Sprintf("number of in-progress v msg: %d", total)
	case "ow":
		total, err := this.getNumberOfWallMsgs()
		if err != nil {
			return err.Error()
		}
		return fmt.Sprintf("total set:ow: %d", total)
	case "pass":
		total, err := this.getNumberOfPassedUsers()
		if err != nil {
			return err.Error()
		}
		return fmt.Sprintf("total set:pass: %d", total)
	}
	return ErrInvalidCmd
}

func handleCommand(openid string, cmdbuf string) string {
	cmds := strings.Fields(cmdbuf)
	if len(cmds) == 0 {
		return ErrInvalidCmd
	}
	var con Console
	if adminCon.IsAdmin(openid) {
		con = adminCon
	} else {
		con = userCon
	}
	var result string

	switch cmds[0] {
	case "h":
		result = con.CmdHelp(cmds[1:])
	case "q":
		fallthrough
	case "?":
		fallthrough
	case "query":
		result = con.CmdQuery(cmds[1:])
	case "s":
		fallthrough
	case "set":
		result = con.CmdSet(cmds[1:])
	default:
		result = ErrInvalidCmd
	}
	return result
}
