package ucrawler

import (
	"WechatWall/crawler/config"
	"testing"
)

func TestParse(t *testing.T) {
	cfg := config.NewForTest()
	ftc := NewFetcher(cfg)
	resp, err := ftc.Do()
	if err != nil {
		t.Fatal(err)
	}
	users, err := Parse(resp)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(users))
	t.Log(users)
}
