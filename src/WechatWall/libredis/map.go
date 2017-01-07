package libredis

import (
	"gopkg.in/redis.v5"

	"time"
)

const (
	USERSMAPNAME = "map:users:openid:"
	PMSGSNAME    = "map:msgs:pending:msgid:"
)

type Map interface {
	Set(string, string) error
	SetTimeout(string, time.Duration) (bool, error)
	Get(string) (string, error)
	Exists(string) (bool, error)
	Del(string) (int64, error)
}

type mMap struct {
	Prefix string
	Client *redis.Client
}

func (this *mMap) Key(k string) string {
	return this.Prefix + k
}

func (this *mMap) Set(k, v string) error {
	key := this.Key(k)
	return this.Client.Set(key, v, 0).Err()
}

func (this *mMap) SetTimeout(k string, timeout time.Duration) (result bool, err error) {
	key := this.Key(k)
	result, err = this.Client.Expire(key, timeout).Result()
	return
}

func (this *mMap) Get(k string) (result string, err error) {
	key := this.Key(k)
	result, err = this.Client.Get(key).Result()
	return
}

// if exists return true
func (this *mMap) Exists(k string) (result bool, err error) {
	key := this.Key(k)
	result, err = this.Client.Exists(key).Result()
	return
}

// How many elements deleted
func (this *mMap) Del(k string) (result int64, err error) {
	key := this.Key(k)
	result, err = this.Client.Del(key).Result()
	return
}

func GetUsersMap() (Map, error) {
	return GetMap(USERSMAPNAME)
}

func GetPMsgsMap() (Map, error) {
	return GetMap(PMSGSNAME)
}

func GetMap(prefix string) (Map, error) {
	client, err := NewRedisClient()
	if err != nil {
		return nil, err
	}
	return &mMap{
		Prefix: prefix,
		Client: client,
	}, nil
}
