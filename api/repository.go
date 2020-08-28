package api

import (
	"database/sql"

	"github.com/giancarlobastos/soccer-manager-api/model"
)

// Repository ...
type Repository struct {
	database *sql.DB
}

// NewRepository ...
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		database: db,
	}
}

func (r *Repository) getPlayer(id int) (model.Player, error) {
	players, err := r.getPlayers("SELECT id, first_name, last_name, age, country, position, market_value, team_id FROM player WHERE id = ?", id)
	return players[0], err
}

func (r *Repository) getAllPlayers() (players []model.Player, err error) {
	return r.getPlayers("SELECT id, first_name, last_name, age, country, position, market_value, team_id FROM player")
}

func (r *Repository) getPlayersByTeamId(TeamId int) (players []model.Player, err error) {
	return r.getPlayers("SELECT id, first_name, last_name, age, country, position, market_value, team_id FROM player WHERE team_id = ?", TeamId)
}

func (r *Repository) getPlayers(query string, args ...interface{}) (players []model.Player, err error) {
	rows, err := r.database.Query(query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var player model.Player
		if err := rows.Scan(
			&player.Id,
			&player.FirstName,
			&player.LastName,
			&player.Age,
			&player.Country,
			&player.Position,
			&player.MarketValue,
			&player.TeamId); err != nil {
			return nil, err
		}
		players = append(players, player)
	}

	return players, nil
}

func (r *Repository) createPlayer(player *model.Player, tx *sql.Tx) error {
	res, err := tx.Exec(
		"INSERT INTO player(first_name, last_name, age, country, position, market_value, team_id) VALUES(?, ?, ?, ?, ?, ?, ?)",
		player.FirstName, player.LastName, player.Age, player.Country, player.Position, player.MarketValue, player.TeamId)

	if err != nil {
		return err
	}

	id, _ := res.LastInsertId()
	player.Id = int(id)
	return nil
}

func (r *Repository) updatePlayer(player *model.Player) error {
	_, err :=
		r.database.Exec("UPDATE player SET first_name = ?, last_name = ?, country = ? WHERE id = ?",
			player.FirstName, player.LastName, player.Country, player.Id)

	return err
}

func (r *Repository) createAccount(account *model.Account, tx *sql.Tx) error {
	res, err := tx.Exec(
		"INSERT INTO account(username, password, first_name, last_name, confirmed, locked, login_attempts, verification_token) VALUES(?, ?, ?, ?, ?, ?, ?, ?)",
		account.Username, account.Password, account.FirstName, account.LastName, account.Confirmed, account.Locked, account.LoginAttempts, account.VerificationToken)

	if err != nil {
		return err
	}

	id, _ := res.LastInsertId()
	account.Id = int(id)
	return nil
}

func (r *Repository) updateAccount(account *model.Account, tx *sql.Tx) error {
	_, err :=
		tx.Exec("UPDATE account SET first_name = ?, last_name = ?, confirmed = ?, locked = ?, login_attempts = ? WHERE id = ?",
			account.FirstName, account.LastName, account.Confirmed, account.Locked, account.LoginAttempts, account.Id)

	return err
}

func (r *Repository) getAccountById(id int) (account model.Account, TeamId int, err error) {
	return r.getAccount(
		"SELECT id, username, password, first_name, last_name, confirmed, locked, login_attempts, verification_token, team_id FROM account WHERE id = ?", id)
}

func (r *Repository) getAccountByUsername(username string) (account model.Account, TeamId int, err error) {
	return r.getAccount(
		"SELECT id, username, password, first_name, last_name, confirmed, locked, login_attempts, verification_token, team_id FROM account WHERE username = ?", username)
}

func (r *Repository) getAccount(query string, args ...interface{}) (account model.Account, TeamId int, err error) {
	return account, TeamId, r.database.QueryRow(query, args...).Scan(
		&account.Id,
		&account.Username,
		&account.Password,
		&account.FirstName,
		&account.LastName,
		&account.Confirmed,
		&account.Locked,
		&account.LoginAttempts,
		&account.VerificationToken,
		&TeamId)
}

func (r *Repository) verifyAccount(token string) error {
	_, err := r.database.Exec("UPDATE account SET confirmed = ? WHERE verification_token = ?", true, token)
	return err
}

func (r *Repository) getTeam(id int) (model.Team, error) {
	team := model.Team{}
	return team, r.database.QueryRow(
		"SELECT id, name, country, available_cash FROM team WHERE id = ?", id).Scan(
		&team.Id,
		&team.Name,
		&team.Country,
		&team.AvailableCash)
}

func (r *Repository) createTeam(team *model.Team, tx *sql.Tx) error {
	res, err := tx.Exec(
		"INSERT INTO team(name, country, available_cash, account_id) VALUES(?, ?, ?, ?)",
		team.Name, team.Country, team.AvailableCash, team.AccountId)

	if err != nil {
		return err
	}

	id, _ := res.LastInsertId()
	team.Id = int(id)
	return nil
}
