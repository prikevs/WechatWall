package libredis

import (
	"encoding/json"
	"time"
)

type Class interface {
	Key() string
	Json() (string, error)
	Loads(string) error
}

type User struct {
	UserOpenid     string `json:"user_openid"`
	UserName       string `json:"user_name"`
	UserCreateTime int64  `json:"user_create_time"`
}

func (this *User) Key() string {
	return this.UserOpenid
}

func (this *User) Json() (string, error) {
	b, err := json.Marshal(this)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (this *User) Loads(j string) error {
	return json.Unmarshal([]byte(j), this)
}

func SetClassToMap(cls Class, mp Map) error {
	data, err := cls.Json()
	if err != nil {
		return err
	}
	return mp.Set(cls.Key(), data)
}

func GetClassFromMap(k string, cls Class, mp Map) error {
	data, err := mp.Get(k)
	if err != nil {
		return err
	}
	if err := cls.Loads(data); err != nil {
		return err
	}
	return nil
}

// message struct
type Msg struct {
	UserOpenid string
	CreateTime int64
	Content    string
	MsgId      string
	MsgType    string
}

func (this *Msg) Key() string {
	return this.MsgId
}

func (this *Msg) Json() (string, error) {
	b, err := json.Marshal(this)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (this *Msg) Loads(j string) error {
	return json.Unmarshal([]byte(j), this)
}

func PublishClassToMQ(cls Class, mq MQ) error {
	data, err := cls.Json()
	if err != nil {
		return err
	}
	return mq.Publish(data)
}

func PublishRClassToMQ(cls Class, mq MQ) error {
	data, err := cls.Json()
	if err != nil {
		return err
	}
	return mq.PublishR(data)
}

func ConsumeClassFromMQ(cls Class, mq MQ, timeout time.Duration) error {
	result, err := mq.Consume(timeout)
	if err != nil {
		return err
	}
	return cls.Loads(result)
}
