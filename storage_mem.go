package main

import "github.com/gofrs/uuid"

type storageMem struct {
	store map[string]string
}

func newStorageMem() storage {
	return &storageMem{
		store: make(map[string]string),
	}
}

func (s storageMem) Create(secret string) (string, error) {
	id := uuid.Must(uuid.NewV4()).String()
	s.store[id] = secret

	return id, nil
}

func (s storageMem) ReadAndDestroy(id string) (string, error) {
	secret, ok := s.store[id]
	if !ok {
		return "", errSecretNotFound
	}

	delete(s.store, id)
	return secret, nil
}
