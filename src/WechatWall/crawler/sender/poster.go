package sender

import (
	"WechatWall/crawler/config"
	"WechatWall/libredis"

	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Poster struct {
	client *http.Client
	cfg    *config.Config
}

func buildReferer(fakeid, token string) string {
	tmpl := "https://mp.weixin.qq.com/cgi-bin/singlesendpage?t=message/send&action=index&tofakeid=%s&token=%s&lang=zh_CN"
	return fmt.Sprintf(tmpl, fakeid, token)
}

func NewPoster(cfg *config.Config) *Poster {
	cfg.SendParams["token"] = cfg.Token
	cfg.SendForm["token"] = cfg.Token
	// to fill: random, content, tofakeid
	return &Poster{
		client: &http.Client{},
		cfg:    cfg,
	}
}

func (this *Poster) Do(msg *libredis.Msg) ([]byte, error) {
	// prepare form data
	newMap := make(map[string]string)
	for k, v := range this.cfg.SendForm {
		newMap[k] = v
	}
	newMap["tofakeid"] = msg.UserOpenid
	newMap["content"] = msg.Content
	newMap["random"] = strconv.FormatFloat(rand.Float64(), 'f', 15, 64)
	form := url.Values{}
	for k, v := range newMap {
		form.Add(k, v)
	}

	req, err := http.NewRequest("POST", this.cfg.SendURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	prepare(this.cfg, req, msg)

	resp, err := this.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return resp_body, nil
}

func prepare(cfg *config.Config, req *http.Request, msg *libredis.Msg) {
	for k, v := range cfg.Headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("Referer", buildReferer(msg.UserOpenid, cfg.Token))
	q := req.URL.Query()
	for k, v := range cfg.SendParams {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
}
