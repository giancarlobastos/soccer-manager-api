package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Router ...
type Router struct {
	service *Service
}

// NewRouter ...
func NewRouter(r *Repository) *Router {
	return &Router{
		service: &Service{
			repository: r,
		},
	}
}

// CreateAccountRequest ...
type CreateAccountRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// Start ...
func (router *Router) Start(addr string) {
	r := mux.NewRouter()
	r.HandleFunc("/accounts", router.createAccount).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", r))
}

func (router *Router) createAccount(w http.ResponseWriter, r *http.Request) {
	var car CreateAccountRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&car); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	account, err := router.service.createAccount(car.FirstName, car.LastName, car.Email, car.Password)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	router.respondCreated(w, r, fmt.Sprintf("/accounts/%d", account.Id))
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (router *Router) respondCreated(w http.ResponseWriter, r *http.Request, path string) {
	w.Header().Set("Path", r.Host+path)
	w.WriteHeader(http.StatusCreated)
}
