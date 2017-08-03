package main

import (
	"fmt"
	"os"

	"github.com/satori/uuid"
	"github.com/xuyu/goredis"
)

type storageRedis struct {
	conn *goredis.Redis
}

func newStorageRedis() (storage, error) {
	if os.Getenv("REDIS_URL") == "" {
		return nil, fmt.Errorf("REDIS_URL environment variable not set")
	}

	c, err := goredis.DialURL(os.Getenv("REDIS_URL"))
	if err != nil {
		return nil, err
	}

	return &storageRedis{
		conn: c,
	}, nil
}

func (s storageRedis) redisKey() string {
	key := os.Getenv("REDIS_KEY")
	if key == "" {
		key = "io.luzifer.ots"
	}

	return key
}

func (s storageRedis) Create(secret string) (string, error) {
	id := uuid.NewV4().String()
	_, err := s.conn.HSet(s.redisKey(), id, secret)

	return id, err
}

func (s storageRedis) ReadAndDestroy(id string) (string, error) {
	secret, err := s.conn.HGet(s.redisKey(), id)
	if err != nil {
		return "", err
	}

	if secret == nil {
		return "", errSecretNotFound
	}

	_, err = s.conn.HDel(s.redisKey(), id)
	return string(secret), err
}
