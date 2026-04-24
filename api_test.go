package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/Luzifer/ots/pkg/customization"
	"github.com/Luzifer/ots/pkg/metrics"
	"github.com/Luzifer/ots/pkg/storage"
	"github.com/Luzifer/ots/pkg/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testCollector = metrics.New()

func TestHandleCreateExpiryOverrideAcceptedValues(t *testing.T) {
	tests := []struct {
		name          string
		expire        int64
		wantExpiresAt bool
	}{
		{
			name:          "zero",
			expire:        0,
			wantExpiresAt: false,
		},
		{
			name:          "one-second",
			expire:        1,
			wantExpiresAt: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			api, _ := newTestAPI(t)
			res := createJSONSecret(api, fmt.Sprintf("/api/create?expire=%d", tc.expire))

			require.Equal(t, http.StatusCreated, res.Code)

			var response apiResponse
			require.NoError(t, json.NewDecoder(res.Body).Decode(&response))
			assert.True(t, response.Success)
			assert.NotEmpty(t, response.SecretID)
			if tc.wantExpiresAt {
				assert.NotNil(t, response.ExpiresAt)
			} else {
				assert.Nil(t, response.ExpiresAt)
			}
		})
	}
}

func TestHandleCreateExpiryOverrideValidation(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		body        string
		expire      string
	}{
		{
			name:        "empty",
			contentType: "application/json",
			body:        `{"secret":"test-secret"}`,
			expire:      "",
		},
		{
			name:        "malformed",
			contentType: "application/json",
			body:        `{"secret":"test-secret"}`,
			expire:      "abc",
		},
		{
			name:        "negative-json",
			contentType: "application/json",
			body:        `{"secret":"test-secret"}`,
			expire:      "-1",
		},
		{
			name:        "negative-form",
			contentType: "application/x-www-form-urlencoded",
			body:        "secret=test-secret",
			expire:      "-1",
		},
		{
			name:        "too-large",
			contentType: "application/json",
			body:        `{"secret":"test-secret"}`,
			expire:      strconv.FormatInt(maxExpirySeconds+1, 10),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			api, store := newTestAPI(t)
			req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/create?expire="+tc.expire, strings.NewReader(tc.body))
			req.Header.Set("Content-Type", tc.contentType)
			res := httptest.NewRecorder()

			api.handleCreate(res, req)

			require.Equal(t, http.StatusBadRequest, res.Code)

			count, err := store.Count()
			require.NoError(t, err)
			assert.Zero(t, count)
		})
	}
}

func createJSONSecret(api *apiServer, target string) *httptest.ResponseRecorder {
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, target, bytes.NewBufferString(`{"secret":"test-secret"}`))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()
	api.handleCreate(res, req)

	return res
}

func newTestAPI(t *testing.T) (*apiServer, storage.Storage) {
	t.Helper()

	oldCfg := cfg
	oldCust := cust
	t.Cleanup(func() {
		cfg = oldCfg
		cust = oldCust
	})

	cfg.SecretExpiry = 3600
	cust = customization.Customize{}

	store := memory.New()
	return newAPI(store, testCollector), store
}
