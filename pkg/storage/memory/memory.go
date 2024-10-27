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
		store           map[string]memStorageSecret
		storePruneTimer *time.Ticker
	}
)

// New creates a new In-Mem storage
func New() storage.Storage {
	store := &storageMem{
		store:           make(map[string]memStorageSecret),
		storePruneTimer: time.NewTicker(time.Minute),
	}

	go store.storePruner()

	return store
}

func (s *storageMem) storePruner() {
	for range s.storePruneTimer.C {
		s.pruneStore()
	}
}

func (s *storageMem) pruneStore() {
	s.Lock()
	defer s.Unlock()

	for k, v := range s.store {
		if v.hasExpired() {
			delete(s.store, k)
		}
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

	// Still check to see if the secret has expired in order to prevent a
	// race condition where a secret has expired but the the store pruner has
	// not yet been invoked.
	if secret.hasExpired() {
		return "", storage.ErrSecretNotFound
	}

	return secret.Secret, nil
}

func (m *memStorageSecret) hasExpired() bool {
	return !m.Expiry.IsZero() && m.Expiry.Before(time.Now())
}
