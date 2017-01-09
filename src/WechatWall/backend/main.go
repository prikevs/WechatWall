package main

import (
	"WechatWall/backend/config"
	"WechatWall/backend/verifier"
	"WechatWall/backend/wall"
	"WechatWall/backend/web"
	"WechatWall/backend/wechat"

	"github.com/gorilla/mux"

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
	r := mux.NewRouter()
	r.HandleFunc("/wx_callback", wechat.CallbackHandler)
	r.HandleFunc("/ws/verifier", verifier.ServeVerifierWS)
	r.HandleFunc("/ws/wall", wall.ServeWallWS)
	r.HandleFunc("/img/{image}", web.ServeImage)
	r.HandleFunc("/static/{type}/{name}", web.ServeStatic)
	r.HandleFunc("/verifier", web.ServeVerifier)
	r.HandleFunc("/wall", web.ServeWall)

	http.Handle("/", r)
}

func main() {
	opts := MustParseArgs()
	cfg := config.New(opts.CfgDir)
	acfg := config.NewAtomicConfig(cfg)

	verifier.Init(acfg)
	wall.Init(acfg)
	wechat.Init(acfg)
	web.Init(acfg)

	log.Println(http.ListenAndServe("127.0.0.1:9999", nil))
}
