package main

import (
	"database/sql"
	"github.com/giancarlobastos/soccer-manager-api/api"
	"github.com/giancarlobastos/soccer-manager-api/repository"
	"github.com/giancarlobastos/soccer-manager-api/service"
	_ "github.com/go-sql-driver/mysql"
	_ "golang.org/x/crypto/bcrypt"
)

var (
	database *sql.DB
	router   *api.Router
)

func main() {
	defer destroy()
	router.Start(":8080")
}

func init() {
	var err error
	database, err = sql.Open("mysql", "root:secret@tcp(mysql:3306)/soccermanager")

	if err != nil {
		panic(err.Error())
	}

	playerRepository := repository.NewPlayerRepository(database)
	teamRepository := repository.NewTeamRepository(database)
	accountRepository := repository.NewAccountRepository(database, teamRepository, playerRepository)
	transferRepository := repository.NewTransferRepository(database)

	accountService := service.NewAccountService(accountRepository, teamRepository, playerRepository)
	playerService := service.NewPlayerService(playerRepository)
	teamService := service.NewTeamService(teamRepository, playerRepository)
	transferService := service.NewTransferService(transferRepository, playerRepository, teamRepository)

	router = api.NewRouter(accountService, teamService, playerService, transferService)
}

func destroy() {
	err := database.Close()

	if err != nil {
		panic(err.Error())
	}
}
