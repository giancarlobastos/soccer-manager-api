package service

import (
	"fmt"
	"github.com/giancarlobastos/soccer-manager-api/domain"
	"github.com/giancarlobastos/soccer-manager-api/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strconv"
)

type AccountService struct {
	accountRepository *repository.AccountRepository
	teamRepository    *repository.TeamRepository
	playerRepository  *repository.PlayerRepository
}

func NewAccountService(
	ar *repository.AccountRepository,
	tr *repository.TeamRepository,
	pr *repository.PlayerRepository) *AccountService {
	return &AccountService{
		accountRepository: ar,
		teamRepository:    tr,
		playerRepository:  pr,
	}
}

func (as *AccountService) GetAccount(accountId int) (*domain.Account, error) {
	account, err := as.accountRepository.GetAccountById(accountId)

	if err != nil {
		return nil, err
	}

	team, err := as.teamRepository.GetTeamByAccountId(accountId)

	if err != nil {
		return nil, err
	}

	account.Team = team

	players, err := as.playerRepository.GetPlayersByTeamId(team.Id)

	if err != nil {
		return nil, err
	}

	account.Team.Players = players

	return &account, nil
}

func (as *AccountService) GetAccountByUsername(username string) (*domain.Account, error) {
	account, err := as.accountRepository.GetAccountByUsername(username)

	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (as *AccountService) CreateAccount(firstName, lastName, email, password string) (*domain.Account, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		panic(err.Error())
	}

	team := as.createTeam(fmt.Sprintf("%s %s's Team", firstName, lastName))
	account := &domain.Account{
		Username:          email,
		Password:          string(encryptedPassword),
		FirstName:         firstName,
		LastName:          lastName,
		VerificationToken: uuid.New().String(),
		LoginAttempts:     0,
		Locked:            false,
		Confirmed:         false,
		Team:              team,
	}

	err = as.accountRepository.CreateAccount(account)

	if err != nil {
		return nil, err
	}

	return account, err
}

func (as *AccountService) createTeam(name string) *domain.Team {
	team := &domain.Team{
		Name:          name,
		Country:       "Brazil",
		AvailableCash: 5000000,
		Players:       make([]domain.Player, 20),
	}

	as.createPlayers(team)
	return team
}

func (as *AccountService) createPlayers(team *domain.Team) {
	as.createPlayersByPosition(team, domain.GoalKeeper, 3)
	as.createPlayersByPosition(team, domain.Defender, 6)
	as.createPlayersByPosition(team, domain.Midfielder, 6)
	as.createPlayersByPosition(team, domain.Forward, 5)
}

func (as *AccountService) createPlayersByPosition(team *domain.Team, position domain.PlayerPosition, quantity int) {
	for i := 0; i < quantity; i++ {
		player := domain.Player{
			FirstName:   string(position),
			LastName:    strconv.Itoa(i),
			Country:     "Brazil",
			Age:         uint8(rand.Intn(22) + 18),
			Position:    position,
			MarketValue: 1000000,
		}
		team.Players = append(team.Players, player)
	}
}
