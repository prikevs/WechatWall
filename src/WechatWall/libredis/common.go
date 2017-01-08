package libredis

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

const (
	MsgTTL = 30 * time.Second
)

type Class interface {
	Key() string
	Json() (string, error)
	Loads(string) error
}

type User struct {
	UserOpenid      string `json:"user_openid"`
	UserName        string `json:"user_name"`
	UserCreateTime  int64  `json:"user_create_time"`
	LastVerifiedMsg string `json:"last_verified_msg"`
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

func SetClassToMapWithTTL(cls Class, mp Map, timeout time.Duration) error {
	// TODO: use transaction here
	if err := SetClassToMap(cls, mp); err != nil {
		return err
	}
	_, err := mp.SetTimeout(cls.Key(), timeout)
	return err
}

func DelClassFromMap(cls Class, mp Map) error {
	affected, err := mp.Del(cls.Key())
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("failed to delete class, affected 0 elements")
	}
	return nil
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
	Username     string
	UserOpenid   string `json:"user_openid"`
	CreateTime   int64  `json:"create_time"`
	Content      string
	MsgId        int64  `json:"msg_id"`
	MsgType      string `json:"msg_type"`
	VerifiedTime int64  `json:"verified_time"`
	TTL          time.Duration
	AddTime      time.Time `json:"add_time"`
}

func (this *Msg) Key() string {
	return strconv.FormatInt(int64(this.MsgId), 10)
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
