package main

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strconv"
)

type Service struct{}

func (s *Service) createAccount(firstName, lastName, email, password string) (*Account, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		panic(err.Error())
	}

	account := &Account{
		Username:          email,
		password:          string(encryptedPassword),
		FirstName:         firstName,
		LastName:          lastName,
		verificationToken: uuid.New().String(),
		loginAttempts:     0,
		locked:            false,
		confirmed:         false,
	}

	var tx *sql.Tx
	tx, err = database.Begin()

	if err != nil {
		panic(err.Error())
	}

	err = repository.createAccount(account, tx)

	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	var team *Team
	team, err = s.createTeam(account, tx)

	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	account.Team = team

	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return account, err
}

func (s *Service) createTeam(account *Account, tx *sql.Tx) (*Team, error) {
	team := &Team{
		Name:          fmt.Sprintf("%s %s's Team", account.FirstName, account.LastName),
		Country:       "Brazil",
		AvailableCash: 5000000,
		Players:       []Player{},
		accountId:     account.Id,
	}

	err := repository.createTeam(team, tx)

	if err != nil {
		return nil, err
	}

	err = s.createPlayers(team, tx)

	if err != nil {
		return nil, err
	}

	return team, err
}

func (s *Service) createPlayers(team *Team, tx *sql.Tx) (err error) {
	if err = s.createPlayersByPosition(team, GoalKeeper, 3, tx); err != nil {
		return
	}

	if err = s.createPlayersByPosition(team, Defender, 6, tx); err != nil {
		return
	}

	if err = s.createPlayersByPosition(team, Midfielder, 6, tx); err != nil {
		return
	}

	if err = s.createPlayersByPosition(team, Forward, 5, tx); err != nil {
		return
	}

	return nil
}

func (s *Service) createPlayersByPosition(team *Team, position PlayerPosition, quantity int, tx *sql.Tx) (err error) {
	for i := 0; i < quantity; i++ {
		player := Player{
			FirstName:   string(position),
			LastName:    strconv.Itoa(i),
			Country:     "Brazil",
			Age:         uint8(rand.Intn(22) + 18),
			Position:    position,
			MarketValue: 1000000,
			teamId:      &team.Id,
		}
		err = repository.createPlayer(&player, tx)

		if err != nil {
			return
		}

		team.Players = append(team.Players, player)
	}
	return nil
}
