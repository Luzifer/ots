package main

import (
	"time"

	"github.com/gofrs/uuid/v3"
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
	id := uuid.Must(uuid.NewV4()).String()
	s.store[id] = memStorageSecret{
		Expiry: time.Now().Add(expireIn),
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
