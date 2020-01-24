package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gofrs/uuid/v3"
	log "github.com/sirupsen/logrus"
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

	s := &storageRedis{
		conn: c,
	}

	if err := s.migrate(); err != nil { // Move from the old to the new storage format
		return nil, err
	}

	return s, nil
}

func (s storageRedis) migrate() error {
	t, err := s.conn.Type(s.redisKey())
	if err != nil {
		return err
	}

	log.Printf("Key %q type: %s", s.redisKey(), t)

	if t == "hash" {
		hashs, err := s.conn.HGetAll(s.redisKey())
		if err != nil {
			return err
		}

		for k, v := range hashs {
			if err := s.writeKey(k, v); err != nil {
				return err
			}
		}

		if _, err = s.conn.Del(s.redisKey()); err != nil {
			return err
		}
	}

	return nil
}

func (s storageRedis) redisExpiry() int {
	var expStr string
	for _, eVar := range []string{"SECRET_EXPIRY", "REDIS_EXPIRY"} {
		if v := os.Getenv(eVar); v != "" {
			expStr = v
			break
		}
	}

	if expStr == "" {
		return 0
	}

	e, err := strconv.ParseInt(expStr, 10, 64)
	if err != nil {
		return 0
	}

	return int(e)
}

func (s storageRedis) redisKey() string {
	key := os.Getenv("REDIS_KEY")
	if key == "" {
		key = "io.luzifer.ots"
	}

	return key
}

func (s storageRedis) Create(secret string) (string, error) {
	id := uuid.Must(uuid.NewV4()).String()
	err := s.writeKey(id, secret)

	return id, err
}

func (s storageRedis) ReadAndDestroy(id string) (string, error) {
	secret, err := s.conn.Get(strings.Join([]string{s.redisKey(), id}, ":"))
	if err != nil {
		return "", err
	}

	if secret == nil {
		return "", errSecretNotFound
	}

	_, err = s.conn.Del(strings.Join([]string{s.redisKey(), id}, ":"))
	return string(secret), err
}

func (s storageRedis) writeKey(id, value string) error {
	return s.conn.Set(
		strings.Join([]string{s.redisKey(), id}, ":"), // Key
		value,           // Secret
		s.redisExpiry(), // Expiry in seconds
		0,               // Expiry milliseconds
		false,           // MustExist
		true,            // MustNotExist
	)
}
