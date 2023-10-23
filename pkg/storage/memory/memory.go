// Package memory implements a pure in-memory store for secrets which
// is suitable for testing and should not be used for productive use
package memory

import (
	"sync"
	"time"

	"github.com/Luzifer/ots/pkg/storage"
	"github.com/gofrs/uuid"
)

type (
	memStorageSecret struct {
		Expiry time.Time
		Secret string
	}

	storageMem struct {
		sync.RWMutex
		store map[string]memStorageSecret
	}
)

// New creates a new In-Mem storage
func New() storage.Storage {
	return &storageMem{
		store: make(map[string]memStorageSecret),
	}
}

func (s *storageMem) Count() (int64, error) {
	s.RLock()
	defer s.RUnlock()

	return int64(len(s.store)), nil
}

func (s *storageMem) Create(secret string, expireIn time.Duration) (string, error) {
	s.Lock()
	defer s.Unlock()

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

func (s *storageMem) ReadAndDestroy(id string) (string, error) {
	s.Lock()
	defer s.Unlock()

	secret, ok := s.store[id]
	if !ok {
		return "", storage.ErrSecretNotFound
	}

	defer delete(s.store, id)

	if !secret.Expiry.IsZero() && secret.Expiry.Before(time.Now()) {
		return "", storage.ErrSecretNotFound
	}

	return secret.Secret, nil
}
