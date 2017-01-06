package ucrawler

import (
	"crawler/config"
	"testing"
)

func TestFetcher(t *testing.T) {
	cfg := config.NewForTest()
	ftc := NewFetcher(cfg)
	if _, err := ftc.Do(); err != nil {
		t.Fatal(err)
	}
}
