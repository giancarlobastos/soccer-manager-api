package repository

import (
	"database/sql"
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

func (tr *TransferRepository) GetTransfers() (transfers []domain.Transfer, err error) {
	query := "SELECT tl.id, tl.asked_price, tl.market_value, " +
		"p.id player_id, p.age, p.country player_country, p.first_name, p.last_name, p.position " +
		"FROM transfer_list tl " +
		"JOIN player p ON p.id = tl.player_id " +
		"WHERE tl.transferred = 0"

	rows, err := tr.db.Query(query)

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
			&transfer.Player.Position); err != nil {
			return nil, err
		}
		transfers = append(transfers, transfer)
	}

	return transfers, nil
}
