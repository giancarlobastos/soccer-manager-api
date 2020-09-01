package repository

import (
	"database/sql"
	"errors"
	"github.com/giancarlobastos/soccer-manager-api/domain"
)

type TransferRepository struct {
	db *sql.DB
}

func NewTransferRepository(db *sql.DB) *TransferRepository {
	return &TransferRepository{
		db: db,
	}
}

func (tr *TransferRepository) NewTransfer(playerId, askedPrice, marketValue int) (transferId int, err error) {
	res, err := tr.db.Exec(
		"INSERT INTO transfer_list(player_id, asked_price, market_value, transferred) VALUES(?, ?, ?, ?)",
		playerId, askedPrice, marketValue, false)

	if err != nil {
		return 0, err
	}

	id, _ := res.LastInsertId()
	return int(id), nil
}

func (tr *TransferRepository) FindTransfers() (transfers []domain.Transfer, err error) {
	query := "SELECT tl.id, tl.asked_price, tl.market_value, " +
		"p.id player_id, p.age, p.country player_country, p.first_name, p.last_name, p.position, p.team_id " +
		"FROM transfer_list tl " +
		"JOIN player p ON p.id = tl.player_id " +
		"WHERE tl.transferred = 0"

	return tr.GetTransfers(query)
}

func (tr *TransferRepository) GetTransfers(query string, args ...interface{}) (transfers []domain.Transfer, err error) {
	rows, err := tr.db.Query(query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var transfer domain.Transfer
		if err := rows.Scan(
			&transfer.Id,
			&transfer.AskedPrice,
			&transfer.MarketValue,
			&transfer.Player.Id,
			&transfer.Player.Age,
			&transfer.Player.Country,
			&transfer.Player.FirstName,
			&transfer.Player.LastName,
			&transfer.Player.Position,
			&transfer.Player.TeamId); err != nil {
			return nil, err
		}
		transfers = append(transfers, transfer)
	}

	return transfers, nil
}

func (tr *TransferRepository) GetTransfer(id int) (domain.Transfer, error) {
	query := "SELECT tl.id, tl.asked_price, tl.market_value, " +
		"p.id player_id, p.age, p.country player_country, p.first_name, p.last_name, p.position, p.team_id " +
		"FROM transfer_list tl " +
		"JOIN player p ON p.id = tl.player_id " +
		"WHERE tl.transferred = 0 AND tl.id = ?"

	transfers, err := tr.GetTransfers(query, id)

	if err != nil {
		return domain.Transfer{}, err
	}

	if len(transfers) == 0 {
		return domain.Transfer{}, sql.ErrNoRows
	}

	return transfers[0], nil
}

func (tr *TransferRepository) ConfirmTransfer(transferId int, buyerId int, newMarketValue int) error {
	tx, err := tr.db.Begin()

	if err != nil {
		return errors.New("internal error")
	}

	res, err :=
		tx.Exec("UPDATE transfer_list tl "+
			"JOIN player p ON p.id = tl.player_id "+
			"JOIN team ts ON ts.id = p.team_id "+
			"JOIN team tb ON tb.id != ts.id "+
			"SET tl.transferred_from = ts.id, tl.transferred_to = tb.id, tl.transferred = 1, "+
			"ts.available_cash = ts.available_cash + tl.asked_price, "+
			"tb.available_cash = tb.available_cash - tl.asked_price, "+
			"p.market_value = ? "+
			"WHERE tl.id = ? AND tl.transferred = 0 AND tb.id = ? AND tb.available_cash >= tl.asked_price",
			newMarketValue, transferId, buyerId)

	if err != nil {
		_ = tx.Rollback()
		return errors.New("transfer not executed")
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected != 4 {
		_ = tx.Rollback()
		return errors.New("transfer not executed")
	}

	return tx.Commit()
}

func (tr *TransferRepository) UpdateTransfer(accountId int, transfer *domain.Transfer) error {
	_, err :=
		tr.db.Exec("UPDATE transfer_list tl "+
			"JOIN player p ON tl.player_id = p.id "+
			"JOIN team t ON p.team_id = t.id "+
			"SET tl.asked_price = ? WHERE tl.id = ? AND t.account_id = ?",
			transfer.AskedPrice, transfer.Id, accountId)

	return err
}
