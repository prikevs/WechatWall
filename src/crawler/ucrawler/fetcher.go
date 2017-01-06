package ucrawler

import (
	"crawler/config"

	"io/ioutil"
	"net/http"
)

type Fetcher struct {
	client *http.Client
	cfg    *config.Config
}

func NewFetcher(cfg *config.Config) *Fetcher {
	return &Fetcher{
		client: &http.Client{},
		cfg:    cfg,
	}
}

func (this *Fetcher) Do() ([]byte, error) {
	req, err := http.NewRequest("GET", this.cfg.UsersURL, nil)
	if err != nil {
		return nil, err
	}
	prepare(this.cfg, req)

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

func prepare(cfg *config.Config, req *http.Request) {
	for k, v := range cfg.Headers {
		req.Header.Set(k, v)
	}
	q := req.URL.Query()
	for k, v := range cfg.UsersParams {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
}
