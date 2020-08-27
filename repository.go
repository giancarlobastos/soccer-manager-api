package main

import "database/sql"

type Repository struct{}

func (r *Repository) getPlayer(id int) (Player, error) {
	players, err := r.getPlayers("SELECT id, first_name, last_name, age, country, position, market_value, team_id FROM player WHERE id = ?", id)

	if err != nil {
		return Player{}, err
	}

	if len(players) == 0 {
		return Player{}, sql.ErrNoRows
	}

	return players[0], nil
}

func (r *Repository) getAllPlayers() (players []Player, err error) {
	return r.getPlayers("SELECT id, first_name, last_name, age, country, position, market_value, team_id FROM player")
}

func (r *Repository) getPlayersByTeamId(teamId int) (players []Player, err error) {
	return r.getPlayers("SELECT id, first_name, last_name, age, country, position, market_value, team_id FROM player WHERE team_id = ?", teamId)
}

func (r *Repository) getPlayers(query string, args ...interface{}) (players []Player, err error) {
	rows, err := database.Query(query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var player Player
		if err := rows.Scan(
			&player.Id,
			&player.FirstName,
			&player.LastName,
			&player.Age,
			&player.Country,
			&player.Position,
			&player.MarketValue,
			&player.teamId); err != nil {
			return nil, err
		}
		players = append(players, player)
	}

	return players, nil
}

func (r *Repository) createPlayer(player *Player, tx *sql.Tx) error {
	res, err := tx.Exec(
		"INSERT INTO player(first_name, last_name, age, country, position, market_value, team_id) VALUES(?, ?, ?, ?, ?, ?, ?)",
		player.FirstName, player.LastName, player.Age, player.Country, player.Position, player.MarketValue, player.teamId)

	if err != nil {
		return err
	}

	id, _ := res.LastInsertId()
	player.Id = int(id)
	return nil
}

func (r *Repository) updatePlayer(player *Player) error {
	_, err :=
		database.Exec("UPDATE player SET first_name = ?, last_name = ?, country = ? WHERE id = ?",
			player.FirstName, player.LastName, player.Country, player.Id)

	return err
}

func (r *Repository) createAccount(account *Account, tx *sql.Tx) error {
	res, err := tx.Exec(
		"INSERT INTO account(username, password, first_name, last_name, confirmed, locked, login_attempts, verification_token) VALUES(?, ?, ?, ?, ?, ?, ?, ?)",
		account.Username, account.password, account.FirstName, account.LastName, account.confirmed, account.locked, account.loginAttempts, account.verificationToken)

	if err != nil {
		return err
	}

	id, _ := res.LastInsertId()
	account.Id = int(id)
	return nil
}

func (r *Repository) updateAccount(account *Account, tx *sql.Tx) error {
	_, err :=
		tx.Exec("UPDATE account SET first_name = ?, last_name = ?, confirmed = ?, locked = ?, login_attempts = ? WHERE id = ?",
			account.FirstName, account.LastName, account.confirmed, account.locked, account.loginAttempts, account.Id)

	return err
}

func (r *Repository) getAccountById(id int) (account Account, err error) {
	return r.getAccount(
		"SELECT id, username, password, first_name, last_name, confirmed, locked, login_attempts, verification_token FROM account WHERE id = ?", id)
}

func (r *Repository) getAccountByUsername(username string) (account Account, err error) {
	return r.getAccount(
		"SELECT id, username, password, first_name, last_name, confirmed, locked, login_attempts, verification_token FROM account WHERE username = ?", username)
}

func (r *Repository) getAccount(query string, args ...interface{}) (account Account, err error) {
	return account, database.QueryRow(query, args...).Scan(
		&account.Id,
		&account.Username,
		&account.password,
		&account.FirstName,
		&account.LastName,
		&account.confirmed,
		&account.locked,
		&account.loginAttempts,
		&account.verificationToken)
}

func (r *Repository) verifyAccount(token string) error {
	_, err := database.Exec("UPDATE account SET confirmed = ? WHERE verification_token = ?", true, token)
	return err
}

func (r *Repository) getTeamById(id int) (Team, error) {
	return r.getTeam("SELECT id, name, country, available_cash FROM team WHERE id = ?", id)
}

func (r *Repository) getTeamByAccountId(id int) (Team, error) {
	return r.getTeam("SELECT id, name, country, available_cash FROM team WHERE account_id = ?", id)
}

func (r *Repository) getTeam(query string, args ...interface{}) (Team, error) {
	team := Team{}
	return team, database.QueryRow(query, args...).Scan(
		&team.Id,
		&team.Name,
		&team.Country,
		&team.AvailableCash)
}

func (r *Repository) createTeam(team *Team, tx *sql.Tx) error {
	res, err := tx.Exec(
		"INSERT INTO team(name, country, available_cash, account_id) VALUES(?, ?, ?, ?)",
		team.Name, team.Country, team.AvailableCash, team.accountId)

	if err != nil {
		return err
	}

	id, _ := res.LastInsertId()
	team.Id = int(id)
	return nil
}
