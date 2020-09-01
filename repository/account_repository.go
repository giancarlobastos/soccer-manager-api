package repository

import (
	"database/sql"
	"github.com/giancarlobastos/soccer-manager-api/domain"
)

type AccountRepository struct {
	db               *sql.DB
	teamRepository   *TeamRepository
	playerRepository *PlayerRepository
}

func NewAccountRepository(
	db *sql.DB,
	tr *TeamRepository,
	pr *PlayerRepository) *AccountRepository {
	return &AccountRepository{
		db:               db,
		teamRepository:   tr,
		playerRepository: pr,
	}
}

func (ar *AccountRepository) CreateAccount(account *domain.Account) error {
	tx, err := ar.db.Begin()

	if err != nil {
		return err
	}

	res, err := tx.Exec(
		"INSERT INTO account(username, password, first_name, last_name, confirmed, locked, login_attempts, verification_token) VALUES(?, ?, ?, ?, ?, ?, ?, ?)",
		account.Username, account.Password, account.FirstName, account.LastName, account.Confirmed, account.Locked, account.LoginAttempts, account.VerificationToken)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	id, err := res.LastInsertId()

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	account.Id = int(id)
	account.Team.AccountId = account.Id

	err = ar.teamRepository.CreateTeam(account.Team, tx)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	for i := 0; i < len(account.Team.Players); i++ {
		account.Team.Players[i].TeamId = &account.Team.Id
		err = ar.playerRepository.CreatePlayer(&account.Team.Players[i], tx)

		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	return err
}

func (ar *AccountRepository) updateAccount(account *domain.Account, tx *sql.Tx) error {
	_, err :=
		tx.Exec("UPDATE account SET first_name = ?, last_name = ?, confirmed = ?, locked = ?, login_attempts = ? WHERE id = ?",
			account.FirstName, account.LastName, account.Confirmed, account.Locked, account.LoginAttempts, account.Id)

	return err
}

func (ar *AccountRepository) GetAccountById(id int) (account domain.Account, err error) {
	return ar.getAccount(
		"SELECT id, username, password, first_name, last_name, confirmed, locked, login_attempts, verification_token FROM account WHERE id = ?", id)
}

func (ar *AccountRepository) GetAccountByUsername(username string) (account domain.Account, err error) {
	return ar.getAccount(
		"SELECT id, username, password, first_name, last_name, confirmed, locked, login_attempts, verification_token FROM account WHERE username = ?", username)
}

func (ar *AccountRepository) getAccount(query string, args ...interface{}) (account domain.Account, err error) {
	return account, ar.db.QueryRow(query, args...).Scan(
		&account.Id,
		&account.Username,
		&account.Password,
		&account.FirstName,
		&account.LastName,
		&account.Confirmed,
		&account.Locked,
		&account.LoginAttempts,
		&account.VerificationToken)
}

func (ar *AccountRepository) VerifyAccount(token string) (bool, error) {
	res, err := ar.db.Exec("UPDATE account SET confirmed = ? WHERE verification_token = ?", true, token)

	if err != nil {
		return false, err
	}

	rowsAffected, _ := res.RowsAffected()

	if rowsAffected != 1 {
		return false, err
	}

	return true, nil
}

func (ar *AccountRepository) RegisterFailedLoginAttempt(username string) error {
	_, err := ar.db.Exec("UPDATE account SET login_attempts = login_attempts + 1, locked = IF(login_attempts > 2, 1, locked) WHERE username = ?", true, username)
	return err
}

func (ar *AccountRepository) ResetLoginAttempts(username string) error {
	_, err := ar.db.Exec("UPDATE account SET login_attempts = 0 WHERE username = ?", true, username)
	return err
}
