package wechat

import (
	"WechatWall/backend/config"
	"WechatWall/libredis"

	"errors"
	"fmt"
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
}

const (
	ErrInvalidCmd    = "invalid command"
	ErrUnauthorized  = "unauthorized usage"
	ErrInternalError = "internal error"
	Success          = "command execute successfully"
)

func (this *AdminConsole) CmdSet(cmds []string) string {
	if len(cmds) == 0 {
		return ErrInvalidCmd
	}
	switch cmds[0] {
	case "needn":
		if len(cmds[1:]) != 1 {
			return ErrInvalidCmd
		}
		switch cmds[1] {
		case "on":
			log.Info("ADMIN set send notification on")
			cfg := config.LoadCfgFromACfg(this.acfg)
			if cfg == nil {
				return ErrInternalError
			}
			cfg.Verifier.SendNotification = true
			this.acfg.StoreConfig(*cfg)
			return Success
		case "off":
			log.Info("ADMIN set send notification off")
			cfg := config.LoadCfgFromACfg(this.acfg)
			if cfg == nil {
				return ErrInternalError
			}
			cfg.Verifier.SendNotification = false
			this.acfg.StoreConfig(*cfg)
			return Success
		default:
			return ErrInvalidCmd
		}
	case "needv":
		if len(cmds[1:]) != 1 {
			return ErrInvalidCmd
		}
		switch cmds[1] {
		case "on":
		case "off":
		default:
			return ErrInvalidCmd
		}

	}

	return ErrInvalidCmd
}

func (this *UserConsole) CmdSet(cmds []string) string {
	return ErrUnauthorized
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
