package main

import (
	"WechatWall/backend/config"
	"WechatWall/backend/lottery"
	"WechatWall/backend/verifier"
	"WechatWall/backend/wall"
	"WechatWall/backend/web"
	"WechatWall/backend/wechat"
	"WechatWall/logger"

	"github.com/goji/httpauth"
	"github.com/gorilla/mux"

	"flag"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

var log = logger.GetLogger("backend")

type Options struct {
	CfgDir string
	TLS    bool
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
	tls := flag.Bool(
		"t",
		false,
		"if use tls")
	flag.Parse()
	opts := &Options{
		CfgDir: *cfgdir,
		TLS:    *tls,
	}
	return opts
}

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/ws/verifier", verifier.ServeVerifierWS)
	r.HandleFunc("/img/{image}", web.ServeImage)
	r.HandleFunc("/static/{type}/{name}", web.ServeStatic)
	r.HandleFunc("/verifier", web.ServeVerifier)
	r.HandleFunc("/wall", web.ServeWall)
	r.HandleFunc("/lottery", lottery.ServeLottery)

	http.HandleFunc("/wx_callback", wechat.CallbackHandler)
	http.HandleFunc("/ws/wall", wall.ServeWallWS)
	http.Handle("/", httpauth.SimpleBasicAuth("kevince", "123456")(r))
	// http.Handle("/", r)
}

func main() {
	opts := MustParseArgs()
	cfg := config.New(opts.CfgDir)
	acfg := config.NewAtomicConfig(cfg)

	verifier.Init(acfg)
	wall.Init(acfg)
	wechat.Init(acfg)
	web.Init(acfg)
	lottery.Init(acfg)

	if opts.TLS {
		crtPath := path.Join(opts.CfgDir, "server.crt")
		keyPath := path.Join(opts.CfgDir, "private.key")
		log.Info(http.ListenAndServeTLS("127.0.0.1:9999", crtPath, keyPath, nil))
	} else {
		log.Info(http.ListenAndServe("127.0.0.1:9999", nil))
	}
}
