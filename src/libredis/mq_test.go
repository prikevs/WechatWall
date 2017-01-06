package libredis

import (
	"testing"
)

func TestMQ(t *testing.T) {
	mq, err := GetMQ("test:mq")
	if err != nil {
		t.Fatal(err)
	}
	if err := mq.Publish("hello"); err != nil {
		t.Fatal(err)
	}
	result, err := mq.Consume(1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}
