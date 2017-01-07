package libredis

import (
	"testing"
)

func TestUserJson(t *testing.T) {
	user := &User{
		UserOpenid:     "123",
		UserName:       "kevince",
		UserCreateTime: 9987,
	}
	js, err := user.Json()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(js)
}

func TestUserLoads(t *testing.T) {
	js := `{"user_openid":"123","user_name":"kevince","user_remark":"","user_create_time":9987}`
	user := &User{}
	if err := user.Loads(js); err != nil {
		t.Fatal(err)
	}
	t.Log(user)
}

func TestSetUserToMap(t *testing.T) {
	user := &User{
		UserOpenid:     "123",
		UserName:       "kevince",
		UserCreateTime: 9987,
	}
	mp, err := GetMap("test:map:2:")
	if err != nil {
		t.Fatal(err)
	}
	if err := SetClassToMap(user, mp); err != nil {
		t.Fatal(err)
	}
}

func TestPublishClassToMQ(t *testing.T) {
	mq, err := GetMQ("test:mq:3")
	if err != nil {
		t.Fatal(err)
	}
	msg := &Msg{
		UserOpenid: "3456",
		CreateTime: 1483750094,
		Content:    "hello",
		MsgId:      6789,
		MsgType:    "text",
	}
	if err := PublishClassToMQ(msg, mq); err != nil {
		t.Fatal("publish, ", err)
	}
	msg2 := &Msg{}
	if err := ConsumeClassFromMQ(msg2, mq, 1); err != nil {
		t.Fatal("consume, ", err)
	}
	t.Log(msg2)
}
