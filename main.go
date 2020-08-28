package main

import (
	"database/sql"

	"github.com/giancarlobastos/soccer-manager-api/api"

	_ "github.com/go-sql-driver/mysql"
	_ "golang.org/x/crypto/bcrypt"
)

var (
	database *sql.DB
)

func main() {
	repository := api.NewRepository(database)
	router := api.NewRouter(repository)

	router.Start(":8080")
	defer destroy()
}

func init() {
	var err error
	database, err = sql.Open("mysql", "root:secret@tcp(127.0.0.1:3306)/soccermanager")

	if err != nil {
		panic(err.Error())
	}
}

func destroy() {
	err := database.Close()

	if err != nil {
		panic(err.Error())
	}
}
