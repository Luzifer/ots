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

type apiResponse struct {
	Success  bool   `json:"success"`
	Error    string `json:"error,omitempty"`
	Secret   string `json:"secret,omitempty"`
	SecretId string `json:"secret_id,omitempty"`
}

type apiRequest struct {
	Secret string `json:"secret"`
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
		tmp := apiRequest{}
		if err := json.NewDecoder(r.Body).Decode(&tmp); err != nil {
			a.errorResponse(res, http.StatusBadRequest, err.Error())
			return
		}
		secret = tmp.Secret
	} else {
		secret = r.FormValue("secret")
	}

	if secret == "" {
		a.errorResponse(res, http.StatusBadRequest, "Secret missing")
		return
	}

	id, err := a.store.Create(secret)
	if err != nil {
		a.errorResponse(res, http.StatusInternalServerError, err.Error())
		return
	}

	a.jsonResponse(res, http.StatusCreated, apiResponse{
		Success:  true,
		SecretId: id,
	})
}

func (a apiServer) handleRead(res http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		a.errorResponse(res, http.StatusBadRequest, "ID missing")
		return
	}

	secret, err := a.store.ReadAndDestroy(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err == errSecretNotFound {
			status = http.StatusNotFound
		}
		a.errorResponse(res, status, err.Error())
		return
	}

	a.jsonResponse(res, http.StatusOK, apiResponse{
		Success: true,
		Secret:  secret,
	})
}

func (a apiServer) errorResponse(res http.ResponseWriter, status int, msg string) {
	a.jsonResponse(res, status, apiResponse{
		Error: msg,
	})
}

func (a apiServer) jsonResponse(res http.ResponseWriter, status int, response apiResponse) {
	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Cache-Control", "no-store, max-age=0")
	res.WriteHeader(status)

	json.NewEncoder(res).Encode(response)
}
