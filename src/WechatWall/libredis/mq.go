package libredis

import (
	"gopkg.in/redis.v5"

	"time"
)

const (
	PMQNAME = "queue:pmq" // MQ for pending messages
	VMQNAME = "queue:vmq" // MQ for verified messages
	SMQNAME = "queue:smq" // MQ for sending messages
)

type MQ interface {
	Length() (int64, error)
	Publish(msg string) error
	PublishR(msg string) error
	Consume(time.Duration) (string, error)
}

type mMQ struct {
	Name   string
	Client *redis.Client
}

func (this *mMQ) Length() (result int64, err error) {
	result, err = this.Client.LLen(this.Name).Result()
	return
}

func (this *mMQ) Publish(msg string) error {
	return this.Client.LPush(this.Name, msg).Err()
}

func (this *mMQ) PublishR(msg string) error {
	return this.Client.RPush(this.Name, msg).Err()
}

func (this *mMQ) Consume(timeout time.Duration) (string, error) {
	result, err := this.Client.BRPop(timeout, this.Name).Result()
	if err != nil {
		return "", err
	}
	return result[1], nil
}

func GetMQ(name string) (MQ, error) {
	client, err := NewRedisClient()
	if err != nil {
		return nil, err
	}
	return &mMQ{
		Name:   name,
		Client: client,
	}, nil
}

func GetPMQ() (MQ, error) {
	return GetMQ(PMQNAME)
}

func GetVMQ() (MQ, error) {
	return GetMQ(VMQNAME)
}

func GetSMQ() (MQ, error) {
	return GetMQ(SMQNAME)
}
