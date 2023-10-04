package main

import (
	"crypto/sha512"
	"encoding/base64"
	"sync"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

var (
	sriCacheStore = newSRICache()
	tplFuncs      template.FuncMap
)

func init() {
	tplFuncs = sprig.FuncMap()
	tplFuncs["assetSRI"] = assetSRIHash
}

func assetSRIHash(assetName string) string {
	if sri, ok := sriCacheStore.Get(assetName); ok {
		return sri
	}

	data, err := assets.ReadFile(assetName)
	if err != nil {
		panic(err)
	}

	h := sha512.New384()
	_, _ = h.Write(data)
	sum := h.Sum(nil)

	sri := "sha384-" + base64.StdEncoding.EncodeToString(sum)
	sriCacheStore.Set(assetName, sri)
	return sri
}

type sriCache struct {
	c map[string]string
	l sync.RWMutex
}

func newSRICache() *sriCache { return &sriCache{c: map[string]string{}} }

func (s *sriCache) Get(assetName string) (string, bool) {
	s.l.RLock()
	defer s.l.RUnlock()

	h, ok := s.c[assetName]
	return h, ok
}

func (s *sriCache) Set(assetName, hash string) {
	s.l.Lock()
	defer s.l.Unlock()

	s.c[assetName] = hash
}
