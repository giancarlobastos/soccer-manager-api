package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "golang.org/x/crypto/bcrypt"
)

var (
	database   *sql.DB
	repository Repository
	service    Service
	router     Router
)

func main() {
	defer destroy()
	router.start(":8080")
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
