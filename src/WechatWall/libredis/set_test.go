package libredis

import (
	"testing"
)

func TestSet(t *testing.T) {
	set, err := GetSet("test:set")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := set.Add("hello", "world", "!"); err != nil {
		t.Fatal(err)
	}
	setb, err := GetSet("test:set:2")
	if _, err := setb.Add("world", "!"); err != nil {
		t.Fatal(err)
	}

	setc, err := set.InterStore(setb, "test:set:3")
	if err != nil {
		t.Fatal(err)
	}
	mems, err := setc.Members()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(mems)
}
