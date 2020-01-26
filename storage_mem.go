package main

import (
	"os"
	"strconv"
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

func (s storageMem) Create(secret string) (string, error) {
	id := uuid.Must(uuid.NewV4()).String()
	s.store[id] = memStorageSecret{
		Expiry: s.expiry(),
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

func (s storageMem) expiry() time.Time {
	exp := os.Getenv("SECRET_EXPIRY")
	if exp == "" {
		return time.Time{}
	}

	e, err := strconv.ParseInt(exp, 10, 64)
	if err != nil {
		return time.Time{}
	}

	return time.Now().Add(time.Duration(e) * time.Second)
}
