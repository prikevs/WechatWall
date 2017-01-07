package filter

import (
	"WechatWall/crawler/config"
	"WechatWall/crawler/ucrawler"
	"WechatWall/crawler/utils"
	"WechatWall/libredis"
	"WechatWall/logger"
)

var log = logger.GetLogger("crawler")

var (
	okset libredis.Set
	wxset libredis.Set
	usmap libredis.Map
)

func gotImage(cfg *config.Config, user *ucrawler.User) bool {
	path := utils.BuildImagePath(cfg, user)
	if utils.FileExists(path) {
		return true
	}
	return false
}

func originCheck(cfg *config.Config, user *ucrawler.User) bool {
	log.Warning("something failed, use origin check")
	if gotImage(cfg, user) {
		return false
	}
	return true
}

func CriticalError(err error) {
	if err != nil {
		log.Critical(err)
	}
}

func FailOnError(err error) {
	if err != nil {
		log.Critical(err)
		panic(err)
	}
}

func init() {
	var err error
	okset, err = libredis.GetOKSet()
	FailOnError(err)
	wxset, err = libredis.GetWXSet()
	FailOnError(err)
	usmap, err = libredis.GetUsersMap()
	FailOnError(err)
}

// return true if need to be crawled
func checkUser(cfg *config.Config, user *ucrawler.User) bool {
	log.Debug("starting to check users")

	// if in ok set
	log.Debug("check if user", user.UserOpenid, "in ok set")
	ismem, err := okset.IsMember(user.UserOpenid)
	if err != nil {
		log.Critical("failed to check if user ok,", err)
		return originCheck(cfg, user)
	}

	if ismem {
		log.Debug("user", user.UserOpenid, "in ok set, skip")
		return false
	}

	// set userinfo to usersmap
	log.Debug("set user", user.UserOpenid, "to usersmap")
	um := &libredis.User{
		UserOpenid:     user.UserOpenid,
		UserName:       user.UserName,
		UserCreateTime: user.UserCreateTime,
	}
	err = libredis.SetClassToMap(um, usmap)
	if err != nil {
		log.Critical("failed to set user to users map,", err)
		return originCheck(cfg, user)
	}

	// delete from wxset
	log.Debug("delete user", user.UserOpenid, "from wxset")
	affected, err := wxset.Del(user.UserOpenid)
	if err != nil {
		log.Critical("failed to del user from wx set,", err)
		return originCheck(cfg, user)
	}

	if affected > 0 {
		log.Info("delete", affected, "elements from wxset")
	} else {
		log.Debug("nothing to do with wxset")
	}

	// check if ready to okset
	if gotImage(cfg, user) {
		// set into okset
		affected, err := okset.Add(user.UserOpenid)
		if err != nil {
			log.Critical("failed to add user to okset, ", err)
			return originCheck(cfg, user)
		}
		if affected > 0 {
			log.Debug("add, success to ADD")
		} else {
			log.Debug("add, already exist")
		}
		return false
	}
	return true

}

func Run(cfg *config.Config, usersch <-chan []ucrawler.User, usersch_filtered chan<- []ucrawler.User) {
	for users := range usersch {
		filtered := make([]ucrawler.User, 0)
		for _, user := range users {
			if checkUser(cfg, &user) {
				filtered = append(filtered, user)
			}
		}
		num := len(filtered)
		if num == 0 {
			log.Info("all users have been fetched, nothing to do")
			continue
		} else {
			log.Info(num, " more images need to be crawled")
		}
		usersch_filtered <- filtered
	}
}
