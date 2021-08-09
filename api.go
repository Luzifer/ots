package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type apiServer struct {
	store storage
}

func newAPI(s storage) *apiServer {
	return &apiServer{
		store: s,
	}
}

func (a apiServer) Register(r *mux.Router) {
	r.HandleFunc("/create", a.handleCreate)
	r.HandleFunc("/get/{id}", a.handleRead)
}

func (a apiServer) handleCreate(res http.ResponseWriter, r *http.Request) {
	var secret string

	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		tmp := map[string]string{}
		if err := json.NewDecoder(r.Body).Decode(&tmp); err != nil {
			a.jsonResponse(res, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		secret = tmp["secret"]
	} else {
		secret = r.FormValue("secret")
	}

	if secret == "" {
		a.jsonResponse(res, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Secret missing",
		})
		return
	}

	id, err := a.store.Create(secret)
	if err != nil {
		a.jsonResponse(res, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	a.jsonResponse(res, http.StatusCreated, map[string]interface{}{
		"success":   true,
		"secret_id": id,
	})
}

func (a apiServer) handleRead(res http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		a.jsonResponse(res, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "ID missing",
		})
		return
	}

	secret, err := a.store.ReadAndDestroy(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err == errSecretNotFound {
			status = http.StatusNotFound
		}
		a.jsonResponse(res, status, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	a.jsonResponse(res, http.StatusOK, map[string]interface{}{
		"success": true,
		"secret":  secret,
	})
}

func (a apiServer) jsonResponse(res http.ResponseWriter, status int, response map[string]interface{}) {
	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Cache-Control", "no-store, max-age=0")
	res.WriteHeader(status)

	json.NewEncoder(res).Encode(response)
}
