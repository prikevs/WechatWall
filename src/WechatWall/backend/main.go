package main

import (
	"WechatWall/backend/config"
	"WechatWall/backend/verifier"
	"WechatWall/backend/wall"
	"WechatWall/backend/wechat"

	"flag"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

type Options struct {
	CfgDir string
}

func getCurrentDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "./"
	}
	return dir
}

func MustParseArgs() *Options {
	cfgdir := flag.String(
		"c",
		path.Join(getCurrentDir(), "etc"),
		"directory of config files, like ./etc")
	flag.Parse()
	opts := &Options{
		CfgDir: *cfgdir,
	}
	return opts
}

func init() {
	http.HandleFunc("/wx_callback", wechat.CallbackHandler)
	http.HandleFunc("/ws/verifier", verifier.ServeVerifierWS)
	http.HandleFunc("/ws/wall", wall.ServeWallWS)
}

func main() {
	opts := MustParseArgs()
	cfg := config.New(opts.CfgDir)

	verifier.Init(&cfg.Verifier)
	wall.Init(&cfg.Wall)
	wechat.Init(&cfg.Wechat)

	log.Println(http.ListenAndServe("127.0.0.1:9999", nil))
}
