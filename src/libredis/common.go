package libredis

import (
	"encoding/json"
)

type User struct {
	UserOpenid     string `json:"user_openid"`
	UserName       string `json:"user_name"`
	UserCreateTime int64  `json:"user_create_time"`
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

func SetUserToMap(user *User, mp Map) error {
	data, err := user.Json()
	if err != nil {
		return err
	}
	return mp.Set(user.UserOpenid, data)
}

func GetUserFromMap(k string, mp Map) (*User, error) {
	user := &User{}
	data, err := mp.Get(k)
	if err != nil {
		return nil, err
	}
	if err := user.Loads(data); err != nil {
		return nil, err
	}
	return user, nil
}

type Msg struct {
}
