package lottery

import (
	"testing"
)

func TestGetLotteryOpenids(t *testing.T) {
	res, err := GetLotteryOpenids(1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}

func TestGetUserinfos(t *testing.T) {
	msgs, err := GetUserInfos()
	if err != nil {
		t.Fatal(err)
	}
	for _, msg := range msgs {
		t.Log(msg)
	}
}
