package api

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/giancarlobastos/soccer-manager-api/model"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Service ...
type Service struct {
	repository *Repository
}

// NewService ...
func NewService(r *Repository) *Service {
	return &Service{
		repository: r,
	}
}

func (s *Service) createAccount(firstName, lastName, email, password string) (*model.Account, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		panic(err.Error())
	}

	account := &model.Account{
		Username:          email,
		Password:          string(encryptedPassword),
		FirstName:         firstName,
		LastName:          lastName,
		VerificationToken: uuid.New().String(),
		LoginAttempts:     0,
		Locked:            false,
		Confirmed:         false,
	}

	var tx *sql.Tx
	tx, err = s.repository.database.Begin()

	if err != nil {
		panic(err.Error())
	}

	err = s.repository.createAccount(account, tx)

	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	var team *model.Team
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

func (s *Service) createTeam(account *model.Account, tx *sql.Tx) (*model.Team, error) {
	team := &model.Team{
		Name:          fmt.Sprintf("%s %s's model.Team", account.FirstName, account.LastName),
		Country:       "Brazil",
		AvailableCash: 5000000,
		Players:       make([]model.Player, 10, 10),
		AccountId:     account.Id,
	}

	err := s.repository.createTeam(team, tx)

	if err != nil {
		return nil, err
	}

	err = s.createPlayers(team, tx)

	if err != nil {
		return nil, err
	}

	return team, err
}

func (s *Service) createPlayers(team *model.Team, tx *sql.Tx) (err error) {
	if err = s.createPlayersByPosition(team, model.GoalKeeper, 3, tx); err != nil {
		return
	}

	if err = s.createPlayersByPosition(team, model.Defender, 6, tx); err != nil {
		return
	}

	if err = s.createPlayersByPosition(team, model.Midfielder, 6, tx); err != nil {
		return
	}

	if err = s.createPlayersByPosition(team, model.Forward, 5, tx); err != nil {
		return
	}

	return nil
}

func (s *Service) createPlayersByPosition(team *model.Team, position model.PlayerPosition, quantity int, tx *sql.Tx) (err error) {
	for i := 0; i < quantity; i++ {
		player := model.Player{
			FirstName:   string(position),
			LastName:    strconv.Itoa(i),
			Country:     "Brazil",
			Age:         uint8(rand.Intn(22) + 18),
			Position:    position,
			MarketValue: 1000000,
			TeamId:      &team.Id,
		}
		err = s.repository.createPlayer(&player, tx)

		if err != nil {
			return
		}

		team.Players = append(team.Players, player)
	}
	return nil
}
