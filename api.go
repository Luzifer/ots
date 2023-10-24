package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Luzifer/ots/pkg/metrics"
	"github.com/Luzifer/ots/pkg/storage"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	errorReasonInvalidJSON    = "invalid_json"
	errorReasonSecretMissing  = "secret_missing"
	errorReasonSecretSize     = "secret_size"
	errorReasonStorageError   = "storage_error"
	errorReasonSecretNotFound = "secret_not_found"
)

type apiServer struct {
	collector *metrics.Collector
	store     storage.Storage
}

type apiResponse struct {
	Success   bool       `json:"success"`
	Error     string     `json:"error,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	Secret    string     `json:"secret,omitempty"`
	SecretID  string     `json:"secret_id,omitempty"`
}

type apiRequest struct {
	Secret string `json:"secret"`
}

func newAPI(s storage.Storage, c *metrics.Collector) *apiServer {
	return &apiServer{
		collector: c,
		store:     s,
	}
}

func (a apiServer) Register(r *mux.Router) {
	r.HandleFunc("/create", a.handleCreate)
	r.HandleFunc("/get/{id}", a.handleRead)
	r.HandleFunc("/isWritable", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) })
	r.HandleFunc("/settings", a.handleSettings).Methods(http.MethodGet)
}

func (a apiServer) handleCreate(res http.ResponseWriter, r *http.Request) {
	if cust.MaxSecretSize > 0 {
		// As a safeguard against HUGE payloads behind a misconfigured
		// proxy we take double the maximum secret size after which we
		// just close the read and cut the connection to the sender.
		r.Body = http.MaxBytesReader(res, r.Body, cust.MaxSecretSize*2) //nolint:gomnd
	}

	var (
		expiry = cfg.SecretExpiry
		secret string
	)

	if !cust.DisableExpiryOverride {
		if ev, err := strconv.ParseInt(r.URL.Query().Get("expire"), 10, 64); err == nil && (ev < expiry || cfg.SecretExpiry == 0) {
			expiry = ev
		}
	}

	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		tmp := apiRequest{}
		if err := json.NewDecoder(r.Body).Decode(&tmp); err != nil {
			if _, ok := err.(*http.MaxBytesError); ok {
				a.collector.CountSecretCreateError(errorReasonSecretSize)
				// We don't do an error response here as the MaxBytesReader
				// automatically cuts the ResponseWriter and we simply cannot
				// answer them.
				return
			}

			a.collector.CountSecretCreateError(errorReasonInvalidJSON)
			a.errorResponse(res, http.StatusBadRequest, err, "decoding request body")
			return
		}
		secret = tmp.Secret
	} else {
		secret = r.FormValue("secret")
	}

	if secret == "" {
		a.collector.CountSecretCreateError(errorReasonSecretMissing)
		a.errorResponse(res, http.StatusBadRequest, errors.New("secret missing"), "")
		return
	}

	if cust.MaxSecretSize > 0 && len(secret) > int(cust.MaxSecretSize) {
		a.collector.CountSecretCreateError(errorReasonSecretSize)
		a.errorResponse(res, http.StatusBadRequest, errors.New("secret size exceeds maximum"), "")
		return
	}

	id, err := a.store.Create(secret, time.Duration(expiry)*time.Second)
	if err != nil {
		a.collector.CountSecretCreateError(errorReasonStorageError)
		a.errorResponse(res, http.StatusInternalServerError, err, "creating secret")
		return
	}

	var expiresAt *time.Time
	if expiry > 0 {
		expiresAt = func(v time.Time) *time.Time { return &v }(time.Now().UTC().Add(time.Duration(expiry) * time.Second))
	}

	a.collector.CountSecretCreated()
	go updateStoredSecretsCount(a.store, a.collector)
	a.jsonResponse(res, http.StatusCreated, apiResponse{
		ExpiresAt: expiresAt,
		Success:   true,
		SecretID:  id,
	})
}

func (a apiServer) handleRead(res http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		a.errorResponse(res, http.StatusBadRequest, errors.New("id missing"), "")
		return
	}

	secret, err := a.store.ReadAndDestroy(id)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, storage.ErrSecretNotFound) {
			a.collector.CountSecretReadError(errorReasonSecretNotFound)
			status = http.StatusNotFound
		} else {
			a.collector.CountSecretReadError(errorReasonStorageError)
		}
		a.errorResponse(res, status, err, "reading & destroying secret")
		return
	}

	a.collector.CountSecretRead()
	go updateStoredSecretsCount(a.store, a.collector)
	a.jsonResponse(res, http.StatusOK, apiResponse{
		Success: true,
		Secret:  secret,
	})
}

func (a apiServer) handleSettings(w http.ResponseWriter, _ *http.Request) {
	a.jsonResponse(w, http.StatusOK, cust)
}

func (a apiServer) errorResponse(res http.ResponseWriter, status int, err error, desc string) {
	errID := uuid.Must(uuid.NewV4()).String()

	if desc != "" {
		// No description: Nothing interesting for the server log
		logrus.WithField("err_id", errID).WithError(err).Error(desc)
	}

	a.jsonResponse(res, status, apiResponse{
		Error: errID,
	})
}

func (apiServer) jsonResponse(res http.ResponseWriter, status int, response any) {
	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Cache-Control", "no-store, max-age=0")
	res.WriteHeader(status)

	if err := json.NewEncoder(res).Encode(response); err != nil {
		logrus.WithError(err).Error("encoding JSON response")
		http.Error(res, `{"error":"could not encode response"}`, http.StatusInternalServerError)
	}
}
