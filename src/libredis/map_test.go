package libredis

import (
	"testing"
)

func TestMAP(t *testing.T) {
	mp, err := GetMap("test:map")
	if err != nil {
		t.Fatal(err)
	}
	if err := mp.Set("123", "read"); err != nil {
		t.Fatal(err)
	}
	exist, err := mp.Exists("123")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(exist)
	result, err := mp.Get("123")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)

	del, err := mp.Del("123")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(del)
	exist, err = mp.Exists("123")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(exist)
}
