package utils

import (
	"WechatWall/crawler/config"
	"WechatWall/crawler/ucrawler"

	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func WriteFile(path string, data []byte) (err error) {
	err = ioutil.WriteFile(path, data, 0644)
	return
}

func BuildImagePath(cfg *config.Config, user *ucrawler.User) string {
	return path.Join(cfg.ImageDir,
		GetFilename(user.UserOpenid, cfg.ImageSuffix))
}

func GetFilename(name, suffix string) string {
	return name + "." + suffix
}

func FileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func CurrentDir() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}
	return dir, nil
}
