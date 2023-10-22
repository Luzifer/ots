// Package redis implements a Redis backed storage for secrets
package redis

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Luzifer/ots/pkg/storage"
	"github.com/gofrs/uuid"
	redis "github.com/redis/go-redis/v9"
)

const (
	redisDefaultPrefix = "io.luzifer.ots"
	redisScanCount     = 10
)

type storageRedis struct {
	conn *redis.Client
}

// New returns a new Redis backed storage
func New() (storage.Storage, error) {
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
		return nil, fmt.Errorf("parsing REDIS_URL: %w", err)
	}

	s := &storageRedis{
		conn: redis.NewClient(opt),
	}

	return s, nil
}

func (s storageRedis) Count() (n int64, err error) {
	var cursor uint64

	for {
		var keys []string

		keys, cursor, err = s.conn.Scan(context.Background(), cursor, s.redisKey("*"), redisScanCount).Result()
		if err != nil {
			return n, fmt.Errorf("scanning stored keys: %w", err)
		}

		n += int64(len(keys))
		if cursor == 0 {
			break
		}
	}

	return n, nil
}

func (s storageRedis) Create(secret string, expireIn time.Duration) (string, error) {
	id := uuid.Must(uuid.NewV4()).String()
	err := s.conn.Set(context.Background(), s.redisKey(id), secret, expireIn).Err()
	if err != nil {
		return "", fmt.Errorf("writing redis key: %w", err)
	}

	return id, nil
}

func (s storageRedis) ReadAndDestroy(id string) (string, error) {
	secret, err := s.conn.Get(context.Background(), s.redisKey(id)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", storage.ErrSecretNotFound
		}
		return "", fmt.Errorf("getting key: %w", err)
	}

	err = s.conn.Del(context.Background(), s.redisKey(id)).Err()
	if err != nil {
		return secret, fmt.Errorf("deleting key: %w", err)
	}
	return secret, nil
}

func (storageRedis) redisKey(id string) string {
	prefix := redisDefaultPrefix
	if prfx := os.Getenv("REDIS_KEY"); prfx != "" {
		prefix = prfx
	}

	return strings.Join([]string{prefix, id}, ":")
}
