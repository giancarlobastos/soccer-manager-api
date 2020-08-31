package service

import (
	"github.com/giancarlobastos/soccer-manager-api/domain"
	"github.com/giancarlobastos/soccer-manager-api/repository"
)

type TransferService struct {
	transferRepository *repository.TransferRepository
	playerRepository   *repository.PlayerRepository
}

func NewTransferService(tr *repository.TransferRepository, pr *repository.PlayerRepository) *TransferService {
	return &TransferService{
		transferRepository: tr,
		playerRepository:   pr,
	}
}

func (ts *TransferService) NewTransfer(playerId, askedPrice int) (transferId int, err error) {
	player, err := ts.playerRepository.GetPlayerOutOfTransferList(playerId)

	if err != nil {
		return 0, err
	}

	transferId, err = ts.transferRepository.NewTransfer(playerId, askedPrice, player.MarketValue)

	if err != nil {
		return 0, err
	}

	return transferId, err
}

func (ts *TransferService) GetTransfers() ([]domain.Transfer, error) {
	transfers, err := ts.transferRepository.GetTransfers()
	return transfers, err
}
