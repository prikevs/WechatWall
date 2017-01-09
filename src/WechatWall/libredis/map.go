package libredis

import (
	"gopkg.in/redis.v5"

	"strings"
	"time"
)

const (
	USERSMAPNAME = "map:users:openid:"
	PMSGSNAME    = "map:msgs:pending:msgid:"
	OWMSGSNAME   = "map:msgs:onwall:msgid:"
	LVMMAPNAME   = "map:users:lvm:openid:"
)

type Map interface {
	Set(string, string) error
	SetTimeout(string, time.Duration) (bool, error)
	Get(string) (string, error)
	Keys() ([]string, error)
	Size() (int64, error)
	Exists(string) (bool, error)
	Del(string) (int64, error)
	TTL(string) (time.Duration, error)
}

type mMap struct {
	Prefix string
	Client *redis.Client
}

func (this *mMap) TTL(k string) (time.Duration, error) {
	key := this.Key(k)
	return this.Client.TTL(key).Result()
}

func (this *mMap) Size() (int64, error) {
	var cursor uint64
	var n int64
	for {
		var keys []string
		var err error
		keys, cursor, err = this.Client.Scan(cursor, this.Prefix+"*", 10).Result()
		if err != nil {
			return 0, err
		}
		n += int64(len(keys))
		if cursor == 0 {
			break
		}
	}

	return n, nil
}

func (this *mMap) Keys() ([]string, error) {
	var cursor uint64
	result := make([]string, 0)
	for {
		var keys []string
		var err error
		keys, cursor, err = this.Client.Scan(cursor, this.Prefix+"*", 10).Result()
		if err != nil {
			return nil, err
		}
		result = append(result, keys...)
		if cursor == 0 {
			break
		}
	}
	for k, v := range result {
		result[k] = strings.TrimPrefix(v, this.Prefix)
	}
	return result, nil
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

// last verifed message map.
// UserOpenid <-> MessageText
func GetLVMMap() (Map, error) {
	return GetMap(LVMMAPNAME)
}

// pending messages map.
// MsgId <-> Msg
func GetPMsgsMap() (Map, error) {
	return GetMap(PMSGSNAME)
}

// already on wall messages map.
// MsgId <-> Msg
func GetOWMsgsMap() (Map, error) {
	return GetMap(OWMSGSNAME)
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
