package repository

import (
	"database/sql"
	"errors"
	"github.com/giancarlobastos/soccer-manager-api/domain"
)

type PlayerRepository struct {
	db *sql.DB
}

func NewPlayerRepository(db *sql.DB) *PlayerRepository {
	return &PlayerRepository{
		db: db,
	}
}

func (pr *PlayerRepository) GetPlayer(id int) (domain.Player, error) {
	players, err := pr.getPlayers("SELECT id, first_name, last_name, age, country, position, market_value, team_id FROM player WHERE id = ?", id)

	if err != nil {
		return domain.Player{}, err
	}

	if len(players) == 0 {
		return domain.Player{}, sql.ErrNoRows
	}

	return players[0], nil
}

func (pr *PlayerRepository) GetPlayerOutOfTransferList(accountId, playerId int) (domain.Player, error) {
	players, err := pr.getPlayers(
		"SELECT p.id, p.first_name, p.last_name, p.age, p.country, p.position, p.market_value, p.team_id "+
			"FROM player p "+
			"JOIN team t ON t.id = p.team_id "+
			"LEFT JOIN transfer_list tl ON p.id = tl.player_id "+
			"WHERE p.id = ? AND t.account_id = ? AND (tl.transferred = 1 OR tl.transferred IS NULL)", playerId, accountId)

	if err != nil {
		return domain.Player{}, err
	}

	if len(players) == 0 {
		return domain.Player{}, errors.New("Invalid player id or it is already in the transfer list")
	}

	return players[0], nil
}

func (pr *PlayerRepository) getAllPlayers() (players []domain.Player, err error) {
	return pr.getPlayers("SELECT id, first_name, last_name, age, country, position, market_value, team_id FROM player")
}

func (pr *PlayerRepository) GetPlayersByTeamId(teamId int) (players []domain.Player, err error) {
	return pr.getPlayers("SELECT id, first_name, last_name, age, country, position, market_value, team_id FROM player WHERE team_id = ?", teamId)
}

func (pr *PlayerRepository) getPlayers(query string, args ...interface{}) (players []domain.Player, err error) {
	rows, err := pr.db.Query(query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var player domain.Player
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

func (pr *PlayerRepository) CreatePlayer(player *domain.Player, tx *sql.Tx) error {
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

func (pr *PlayerRepository) UpdatePlayer(accountId int, player *domain.Player) error {
	_, err :=
		pr.db.Exec("UPDATE player p "+
			"JOIN team t ON t.id = p.team_id "+
			"SET p.first_name = ?, p.last_name = ?, p.country = ? "+
			"WHERE p.id = ? AND t.account_id = ?",
			player.FirstName, player.LastName, player.Country, player.Id, accountId)

	return err
}
