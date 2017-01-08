package libredis

import (
	"gopkg.in/redis.v5"
)

const (
	OKSETNAME      = "set:ok"
	OWSETNAME      = "set:ow" // messages on wall set
	WXSETNAME      = "set:wx"
	PASSSETNAME    = "set:pass"
	SENTSETNAME    = "set:sent"
	LOTTERYSETNAME = "set:lottery"
)

type Set interface {
	Total() (int64, error)
	GetSetName() string
	Members() ([]string, error)
	IsMember(string) (bool, error)
	Add(...string) (int64, error)
	Del(...string) (int64, error)
	InterStore(Set, string) (Set, error)
}

type mSet struct {
	Name   string
	Client *redis.Client
}

func (this *mSet) Total() (int64, error) {
	return this.Client.SCard(this.Name).Result()
}

func (this *mSet) GetSetName() string {
	return this.Name
}

func (this *mSet) IsMember(v string) (result bool, err error) {
	result = false
	result, err = this.Client.SIsMember(this.Name, v).Result()
	return
}

func (this *mSet) Members() (result []string, err error) {
	result, err = this.Client.SMembers(this.Name).Result()
	return
}

func (this *mSet) Add(vs ...string) (result int64, err error) {
	new := make([]interface{}, len(vs))
	for i, v := range vs {
		new[i] = v
	}
	result, err = this.Client.SAdd(this.Name, new...).Result()
	return
}

func (this *mSet) Del(vs ...string) (result int64, err error) {
	new := make([]interface{}, len(vs))
	for i, v := range vs {
		new[i] = v
	}
	result, err = this.Client.SRem(this.Name, new...).Result()
	return
}

func (this *mSet) InterStore(set Set, target string) (Set, error) {
	_, err := this.Client.SInterStore(target, set.GetSetName()).Result()
	if err != nil {
		return nil, err
	}
	result, err := GetSet(target)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetSet(name string) (Set, error) {
	client, err := NewRedisClient()
	if err != nil {
		return nil, err
	}
	return &mSet{
		Name:   name,
		Client: client,
	}, nil
}

func GetOKSet() (Set, error) {
	return GetSet(OKSETNAME)
}

func GetPassSet() (Set, error) {
	return GetSet(PASSSETNAME)
}

func GetWXSet() (Set, error) {
	return GetSet(WXSETNAME)
}

func GetSentSet() (Set, error) {
	return GetSet(SENTSETNAME)
}

func GetOWSet() (Set, error) {
	return GetSet(OWSETNAME)
}
