package icrawler

import (
	"WechatWall/crawler/config"
	"WechatWall/crawler/ucrawler"
	"testing"
)

func TestIFetcher(t *testing.T) {
	cfg := config.NewForTest()
	user := &ucrawler.User{
		UserOpenid: "o_5d9syXDq7AiyB387ZHEpX6NKQE",
	}
	ftc := NewIFetcher(cfg, user)

	resp, err := ftc.Do()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}
