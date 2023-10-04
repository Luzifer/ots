package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	redis "github.com/redis/go-redis/v9"
)

const redisDefaultPrefix = "io.luzifer.ots"

type storageRedis struct {
	conn *redis.Client
}

func newStorageRedis() (storage, error) {
	if os.Getenv("REDIS_URL") == "" {
		return nil, fmt.Errorf("REDIS_URL environment variable not set")
	}

	// We replace the old URI format
	//		tcp://auth:password@127.0.0.1:6379/0
	// with the new one
	//		redis://<user>:<password>@<host>:<port>/<db_number>
	// in order to maintain backwards compatibility
	opt, err := redis.ParseURL(strings.Replace(os.Getenv("REDIS_URL"), "tcp://", "redis://", 1))
	if err != nil {
		return nil, errors.Wrap(err, "parsing REDIS_URL")
	}

	s := &storageRedis{
		conn: redis.NewClient(opt),
	}

	return s, nil
}

func (s storageRedis) Create(secret string, expireIn time.Duration) (string, error) {
	id := uuid.Must(uuid.NewV4()).String()
	err := s.conn.Set(context.Background(), s.redisKey(id), secret, expireIn).Err()

	return id, errors.Wrap(err, "writing redis key")
}

func (s storageRedis) ReadAndDestroy(id string) (string, error) {
	secret, err := s.conn.Get(context.Background(), s.redisKey(id)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", errSecretNotFound
		}
		return "", errors.Wrap(err, "getting key")
	}

	err = s.conn.Del(context.Background(), s.redisKey(id)).Err()
	return secret, errors.Wrap(err, "deleting key")
}

func (storageRedis) redisKey(id string) string {
	prefix := redisDefaultPrefix
	if prfx := os.Getenv("REDIS_KEY"); prfx != "" {
		prefix = prfx
	}

	return strings.Join([]string{prefix, id}, ":")
}
