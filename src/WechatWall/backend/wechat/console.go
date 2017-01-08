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
	vMQ     libredis.MQ
	pMQ     libredis.MQ
	owSet   libredis.Set
	passSet libredis.Set
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
	return this.owSet.Total()
}
func (this *UserConsole) getNumberOfPassedUsers() (num int64, err error) {
	return this.passSet.Total()
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

	userCon := &UserConsole{
		vMQ:     vMQ,
		pMQ:     pMQ,
		owSet:   owSet,
		passSet: passSet,
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

func (this *AdminConsole) CmdSet(cmds []string) string {
	if len(cmds) == 0 {
		return ErrInvalidCmd
	}
	switch cmds[0] {
	case "sendn":
		if len(cmds[1:]) != 1 {
			return ErrInvalidCmd
		}
		switch cmds[1] {
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
	case "needv":
		if len(cmds[1:]) != 1 {
			return ErrInvalidCmd
		}
		switch cmds[1] {
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
	case "ttl":
		if len(cmds[1:]) != 1 {
			return ErrInvalidCmd
		}
		ttl, err := strconv.Atoi(cmds[1])
		if err != nil {
			return ErrInvalidCmd
		}
		return SetConfig(this.acfg, ttl,
			func(cfg *config.Config, val interface{}) {
				cfg.Verifier.MaxMsgWaitingTime = val.(int)
			})

	case "svd":

	}

	return ErrInvalidCmd
}

func (this *UserConsole) CmdSet(cmds []string) string {
	return ErrUnauthorized
}

func (this *AdminConsole) CmdHelp(cmds []string) string {
	return this.UserConsole.CmdHelp(cmds) + `Admin Console Help:
query/q/?
ttl:
    current ttl of unverified message
svd:
    current duration of sending message to verification
sendn:
    current status of if sending notification
needv:
    current status of if needing verification

set/s
ttl <int>:
    set ttl of unverified message (in seconds)
svd <int>:
    set duration of sending message to verification (in seconds)
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
		case "ttl":
			ttl := GetConfig(this.acfg,
				func(cfg *config.Config) interface{} {
					return cfg.Verifier.MaxMsgWaitingTime
				}).(int)
			result = fmt.Sprintf("current TTL is: %d", ttl)
		case "svd":
		case "sendn":
		case "needv":
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
