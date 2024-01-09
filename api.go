package main

import (
	cache2 "WB/cache"
	"WB/storage"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type APIServer struct {
	listenAddr string
	store      storage.Storage
	stream     Stream
	cache      cache2.Cache
}

func NewAPIServer(listerAddr string, store storage.Storage) *APIServer {
	return &APIServer{
		listenAddr: listerAddr,
		store:      store,
	}
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteSegment(w, r)
	}
	return fmt.Errorf("method now allowed %s", r.Method)

}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	uid := mux.Vars(r)["id"]
	user, err := s.cache.Order().Find(uid)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, user)
}

func (s *APIServer) handleDeleteSegment(w http.ResponseWriter, r *http.Request) error {
	segment := mux.Vars(r)["slug"]
	if err := s.store.DeleteSegment(segment); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, "segment deleted sucsessfully")
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func (s *APIServer) Run() {

	router := mux.NewRouter()
	router.HandleFunc("/user/{id}", makeHTTPHandleFunc(s.handleGetAccount))
	router.HandleFunc("/user/{slug}", makeHTTPHandleFunc(s.handleDeleteSegment))
	http.ListenAndServe(s.listenAddr, router)

	log.Panicln("API server running on port", s.listenAddr)

}

type ApiError struct {
	Error string `json:"error"`
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}
