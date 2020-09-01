package api

import (
	"encoding/json"
	"fmt"
	"github.com/giancarlobastos/soccer-manager-api/security"
	"github.com/giancarlobastos/soccer-manager-api/service"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Router struct {
	accountService           *service.AccountService
	teamService              *service.TeamService
	playerService            *service.PlayerService
	transferService          *service.TransferService
	authenticationMiddleware *security.AuthenticationMiddleware
}

func NewRouter(as *service.AccountService, ts *service.TeamService, ps *service.PlayerService, tfs *service.TransferService) *Router {
	amw := security.NewAuthenticationMiddleware(as,
		map[string]string{
			"getAccount":      "USER",
			"getPlayer":       "USER",
			"updatePlayer":    "USER",
			"getTeam":         "USER",
			"updateTeam":      "USER",
			"newTransfer":     "USER",
			"getTransfers":    "USER",
			"confirmTransfer": "USER",
			"updateTransfer":  "USER",
		})
	return &Router{
		accountService:           as,
		teamService:              ts,
		playerService:            ps,
		transferService:          tfs,
		authenticationMiddleware: amw,
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
	r.HandleFunc("/accounts/{accountId}", router.getAccount).Methods("GET").Name("getAccount")
	r.HandleFunc("/verify-account", router.verifyAccount).Queries("token", "{token}").Methods("GET")
	r.HandleFunc("/players/{playerId}", router.getPlayer).Methods("GET").Name("getPlayer")
	r.HandleFunc("/players/{playerId}", router.updatePlayer).Methods("PATCH").Name("updatePlayer")
	r.HandleFunc("/teams/{teamId}", router.getTeam).Methods("GET").Name("getTeam")
	r.HandleFunc("/teams/{teamId}", router.updateTeam).Methods("PATCH").Name("updateTeam")
	r.HandleFunc("/transfers", router.newTransfer).Methods("POST").Name("newTransfer")
	r.HandleFunc("/transfers", router.getTransfers).Methods("GET").Name("getTransfers")
	r.HandleFunc("/transfers/{transferId}", router.confirmTransfer).Methods("PUT").Name("confirmTransfer")
	r.HandleFunc("/transfers/{transferId}", router.updateTransfer).Methods("PATCH").Name("updateTransfer")

	r.HandleFunc("/authenticate", router.authenticationMiddleware.Authenticate).Methods("POST")
	r.Use(router.authenticationMiddleware.Middleware)
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

	principal := router.authenticationMiddleware.GetPrincipal(r)

	if principal.AccountId != accountId {
		respondWithError(w, http.StatusNotFound, "account id not found")
		return
	}

	account, err := router.accountService.GetAccount(accountId)

	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, account)
}

func (router *Router) verifyAccount(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	if len(token) > 0 && router.accountService.VerifyAccount(token) {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	respondWithError(w, http.StatusBadRequest, "Invalid token")
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

	principal := router.authenticationMiddleware.GetPrincipal(r)
	player, err := router.playerService.UpdatePlayer(principal.AccountId, playerId, patchJSON)

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

	principal := router.authenticationMiddleware.GetPrincipal(r)
	team, err := router.teamService.GetTeam(principal.AccountId, teamId)

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

	principal := router.authenticationMiddleware.GetPrincipal(r)
	team, err := router.teamService.UpdateTeam(principal.AccountId, teamId, patchJSON)

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

	principal := router.authenticationMiddleware.GetPrincipal(r)
	transferId, err := router.transferService.NewTransfer(principal.AccountId, tr.PlayerId, tr.AskedPrice)

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

func (router *Router) confirmTransfer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transferId, err := strconv.Atoi(vars["transferId"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	principal := router.authenticationMiddleware.GetPrincipal(r)
	err = router.transferService.ConfirmTransfer(principal.AccountId, transferId)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to confirm transfer")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (router *Router) updateTransfer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transferId, err := strconv.Atoi(vars["transferId"])

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

	principal := router.authenticationMiddleware.GetPrincipal(r)
	team, err := router.transferService.UpdateTransfer(principal.AccountId, transferId, patchJSON)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Error in updating team's data")
		return
	}

	respondWithJSON(w, http.StatusOK, team)
}
