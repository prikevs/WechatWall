package icrawler

import (
	"WechatWall/crawler/config"
	"WechatWall/crawler/ucrawler"

	"errors"
	// "fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type IFetcher struct {
	client *http.Client
	cfg    *config.Config
	user   *ucrawler.User
}

func NewIFetcher(cfg *config.Config, user *ucrawler.User) *IFetcher {
	return &IFetcher{
		client: &http.Client{},
		cfg:    cfg,
		user:   user,
	}
}

func (this *IFetcher) Do() ([]byte, error) {
	req, err := http.NewRequest("GET", this.cfg.ImageURL, nil)
	if err != nil {
		return nil, err
	}
	prepare(req, this.cfg, this.user)

	resp, err := this.client.Do(req)
	if err != nil {
		return nil, err
	}
	v, ok := resp.Header["Content-Type"]
	if !ok || !strings.Contains(v[0], "image") {
		return nil, errors.New("invalid type")
	}
	defer resp.Body.Close()

	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return resp_body, nil
}

func prepare(req *http.Request, cfg *config.Config, user *ucrawler.User) {
	for k, v := range cfg.Headers {
		req.Header.Set(k, v)
	}
	q := req.URL.Query()

	newMap := make(map[string]string)
	for k, v := range cfg.ImageParams {
		newMap[k] = v
	}
	newMap["fakeid"] = user.UserOpenid
	for k, v := range newMap {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
}
