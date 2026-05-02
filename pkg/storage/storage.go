// Package storage describes the requirements a storage provider
// has to fulfill ot be usable in OTS
package storage

import (
	"errors"
	"time"
)

type (
	// Storage is the interface to implement in each storage provider
	Storage interface {
		// Count returns the number of stored secrets
		Count() (int64, error)
		// Create inserts a new secret and returns its ID
		Create(secret string, expireIn time.Duration) (string, error)
		// ReadAndDestroy returns a secret and while reading removes it
		// from the storage
		ReadAndDestroy(id string) (string, error)
	}
)

// ErrSecretNotFound is a generic error to be returned when a secret
// does not exist in the backend. It will then be handled by API.
var ErrSecretNotFound = errors.New("secret not found")
