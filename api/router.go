package api

import (
	"encoding/json"
	"fmt"
	"github.com/giancarlobastos/soccer-manager-api/service"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Router struct {
	accountService  *service.AccountService
	teamService     *service.TeamService
	playerService   *service.PlayerService
	transferService *service.TransferService
}

func NewRouter(as *service.AccountService, ts *service.TeamService, ps *service.PlayerService, tfs *service.TransferService) *Router {
	return &Router{
		accountService:  as,
		teamService:     ts,
		playerService:   ps,
		transferService: tfs,
	}
}

type CreateAccountRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type PutInTransferListRequest struct {
	PlayerId   int
	AskedPrice int
}

func (router *Router) Start(addr string) {
	r := mux.NewRouter()

	r.HandleFunc("/accounts", router.createAccount).Methods("POST")
	r.HandleFunc("/accounts/{accountId}", router.getAccount).Methods("GET")
	r.HandleFunc("/players/{playerId}", router.getPlayer).Methods("GET")
	r.HandleFunc("/players/{playerId}", router.updatePlayer).Methods("PATCH")
	r.HandleFunc("/teams/{teamId}", router.getTeam).Methods("GET")
	r.HandleFunc("/teams/{teamId}", router.updateTeam).Methods("PATCH")
	r.HandleFunc("/transfers", router.newTransfer).Methods("POST")
	r.HandleFunc("/transfers", router.getTransfers).Methods("GET")

	r.HandleFunc("/authenticate", router.createAccount).Methods("POST")
	r.HandleFunc("/verify-account", router.createAccount).Methods("GET")
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

	account, err := router.accountService.CreateAccount(car.FirstName, car.LastName, car.Email, car.Password)

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

	account, err := router.accountService.GetAccount(accountId)

	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, account)
}

func (router *Router) getPlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playerId, err := strconv.Atoi(vars["playerId"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	player, err := router.playerService.GetPlayer(playerId)

	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, player)
}

func (router *Router) updatePlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playerId, err := strconv.Atoi(vars["playerId"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	patchJSON, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	player, err := router.playerService.UpdatePlayer(playerId, patchJSON)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error in updating player's data")
		return
	}

	respondWithJSON(w, http.StatusOK, player)
}

func (router *Router) getTeam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamId, err := strconv.Atoi(vars["teamId"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	team, err := router.teamService.GetTeam(teamId)

	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, team)
}

func (router *Router) updateTeam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamId, err := strconv.Atoi(vars["teamId"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	patchJSON, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	team, err := router.teamService.UpdateTeam(teamId, patchJSON)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Error in updating team's data")
		return
	}

	respondWithJSON(w, http.StatusOK, team)
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

func (router *Router) newTransfer(w http.ResponseWriter, r *http.Request) {
	var tr PutInTransferListRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&tr); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	transferId, err := router.transferService.NewTransfer(tr.PlayerId, tr.AskedPrice)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	router.respondCreated(w, r, fmt.Sprintf("/transfers/%d", transferId))
}

func (router *Router) getTransfers(w http.ResponseWriter, r *http.Request) {
	transfers, err := router.transferService.GetTransfers()

	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, transfers)
}
