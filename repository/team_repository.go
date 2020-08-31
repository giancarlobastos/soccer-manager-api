package repository

import (
	"database/sql"
	"github.com/giancarlobastos/soccer-manager-api/domain"
)

type TeamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) *TeamRepository {
	return &TeamRepository{
		db: db,
	}
}

func (tr *TeamRepository) GetTeamById(id int) (*domain.Team, error) {
	return tr.getTeam("SELECT id, name, country, available_cash FROM team WHERE id = ?", id)
}

func (tr *TeamRepository) GetTeamByAccountId(id int) (*domain.Team, error) {
	return tr.getTeam("SELECT id, name, country, available_cash FROM team WHERE account_id = ?", id)
}

func (tr *TeamRepository) getTeam(query string, args ...interface{}) (*domain.Team, error) {
	team := &domain.Team{}
	return team, tr.db.QueryRow(query, args...).Scan(
		&team.Id,
		&team.Name,
		&team.Country,
		&team.AvailableCash)
}

func (tr *TeamRepository) CreateTeam(team *domain.Team, tx *sql.Tx) error {
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

func (tr *TeamRepository) UpdateTeam(team *domain.Team) error {
	_, err :=
		tr.db.Exec("UPDATE team SET name = ?, country = ? WHERE id = ?", team.Name, team.Country, team.Id)

	return err
}
