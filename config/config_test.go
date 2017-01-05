package config

import (
	"testing"
)

func TestNew(t *testing.T) {
	c := New()
	t.Log(c)
}
