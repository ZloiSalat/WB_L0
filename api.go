package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sync"
)

type APIServer struct {
	listenAddr string
	store      Storage
	stream     Stream
}

func NewAPIServer(listerAddr string, store Storage) *APIServer {
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
	user, err := s.store.GetActiveUsers()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, user)
}

func (s *APIServer) handleCreateUser(wg *sync.WaitGroup, notifyChannel chan []byte) error {
	defer wg.Done()

	for {
		// Read from the notifyChannel
		data, ok := <-notifyChannel
		if !ok {
			// Channel closed, exiting the loop
			return nil
		}

		createUserReq := new(User)
		if err := json.Unmarshal(data, createUserReq); err != nil {
			return err
		}
		user, _ := NewUser(
			createUserReq.OrderUid,
			createUserReq.TrackNumber,
			createUserReq.Entry,
			createUserReq.Delivery,
			createUserReq.Payment,
			createUserReq.Items,
			createUserReq.Locale,
			createUserReq.InternalSignature,
			createUserReq.CustomerID,
			createUserReq.DeliveryService,
			createUserReq.Shardkey,
			createUserReq.SmID,
			createUserReq.DateCreated,
			createUserReq.OofShard,
		)
		if err := s.store.CreateUser(user); err != nil {
			return err
		}
	}

}

func (s *APIServer) handleDeleteSegment(w http.ResponseWriter, r *http.Request) error {
	segment := mux.Vars(r)["slug"]
	if err := s.store.DeleteSegment(segment); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, "segment deleted sucsessfully")
}

/*func (s *APIServer) handleAddUserToSegment(w http.ResponseWriter, r *http.Request) error {
	createSegmentRequest := new(Request)
	if err := json.NewDecoder(r.Body).Decode(createSegmentRequest); err != nil {
		return err
	}
	segment := NewRequest(createSegmentRequest.UserID, createSegmentRequest.AddSegments, createSegmentRequest.RemoveSegments)
	if err := s.store.AddUserToSegment(segment); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, segment)
}*/

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
