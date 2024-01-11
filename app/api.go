package app

import (
	"WB/cache"
	"WB/store"
	"WB/types"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type APIServer struct {
	listenAddr string
	store      store.Store
	cache      cache.Cache
}

func NewAPIServer(listerAddr string, store store.Store, cache cache.Cache) *APIServer {
	return &APIServer{
		listenAddr: listerAddr,
		store:      store,
		cache:      cache,
	}
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}

	return fmt.Errorf("method now allowed %s", r.Method)

}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	var user types.User

	uid := mux.Vars(r)["id"]
	order, err := s.cache.Order().Find(uid)
	if err != nil {
		return err
	}

	err = json.Unmarshal(order.Data, &user)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, user)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func (s *APIServer) Run() {

	router := mux.NewRouter()
	router.HandleFunc("/user/{id}", makeHTTPHandleFunc(s.handleGetAccount))
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
