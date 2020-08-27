package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Router struct{}

type CreateAccountRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (router *Router) start(addr string) {
	r := mux.NewRouter()
	r.HandleFunc("/authenticate", router.createAccount).Methods("POST")

	r.HandleFunc("/accounts", router.createAccount).Methods("POST")
	r.HandleFunc("/accounts/{accountId}", router.getAccount).Methods("GET")

	r.HandleFunc("/verify-account", router.createAccount).Methods("GET")
	r.HandleFunc("/players/{playerId}", router.createAccount).Methods("GET")
	r.HandleFunc("/players/{playerId}", router.createAccount).Methods("PATCH")
	r.HandleFunc("/teams/{teamId}", router.createAccount).Methods("GET")
	r.HandleFunc("/teams/{teamId}", router.createAccount).Methods("PATCH")
	r.HandleFunc("/transfers", router.createAccount).Methods("GET")
	r.HandleFunc("/transfers", router.createAccount).Methods("POST")
	r.HandleFunc("/transfers/{transferId}", router.createAccount).Methods("PATCH")
	r.HandleFunc("/transfers/{transferId}", router.createAccount).Methods("PUT")

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

	account, err := service.createAccount(car.FirstName, car.LastName, car.Email, car.Password)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	router.respondCreated(w, r, fmt.Sprintf("/accounts/%d", account.Id))
}

func (router *Router) getAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountId, err := strconv.Atoi(vars["accountId"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	account, err := service.getAccount(accountId)

	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, account)
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
