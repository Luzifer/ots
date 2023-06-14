package main

import (
	"time"

	"github.com/gofrs/uuid"
)

type memStorageSecret struct {
	Expiry time.Time
	Secret string
}

type storageMem struct {
	store map[string]memStorageSecret
}

func newStorageMem() storage {
	return &storageMem{
		store: make(map[string]memStorageSecret),
	}
}

func (s storageMem) Create(secret string, expireIn time.Duration) (string, error) {
	var (
		expire time.Time
		id     = uuid.Must(uuid.NewV4()).String()
	)

	if expireIn > 0 {
		expire = time.Now().Add(expireIn)
	}

	s.store[id] = memStorageSecret{
		Expiry: expire,
		Secret: secret,
	}

	return id, nil
}

func (s storageMem) ReadAndDestroy(id string) (string, error) {
	secret, ok := s.store[id]
	if !ok {
		return "", errSecretNotFound
	}

	defer delete(s.store, id)

	if !secret.Expiry.IsZero() && secret.Expiry.Before(time.Now()) {
		return "", errSecretNotFound
	}

	return secret.Secret, nil
}
