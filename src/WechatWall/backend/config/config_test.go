package config

import (
	"testing"
)

func TestNew(t *testing.T) {
	dir := "../etc"
	cfg := New(dir)
	t.Log(cfg)
}
