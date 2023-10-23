package main

import (
	"fmt"

	"github.com/Luzifer/ots/pkg/storage"
	"github.com/Luzifer/ots/pkg/storage/memory"
	"github.com/Luzifer/ots/pkg/storage/redis"
)

func getStorageByType(t string) (storage.Storage, error) {
	switch t {
	case "mem":
		return memory.New(), nil

	case "redis":
		s, err := redis.New()
		if err != nil {
			return s, fmt.Errorf("creating redis storage: %w", err)
		}
		return s, nil

	default:
		return nil, fmt.Errorf("storage type %q not found", t)
	}
}
