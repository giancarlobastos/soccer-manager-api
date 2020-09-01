package service

import (
	"encoding/json"
	"errors"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/giancarlobastos/soccer-manager-api/domain"
	"github.com/giancarlobastos/soccer-manager-api/repository"
	"math/rand"
)

type TransferService struct {
	transferRepository *repository.TransferRepository
	playerRepository   *repository.PlayerRepository
	teamRepository     *repository.TeamRepository
}

func NewTransferService(tfr *repository.TransferRepository, pr *repository.PlayerRepository, tr *repository.TeamRepository) *TransferService {
	return &TransferService{
		transferRepository: tfr,
		playerRepository:   pr,
		teamRepository:     tr,
	}
}

func (ts *TransferService) NewTransfer(accountId, playerId, askedPrice int) (transferId int, err error) {
	player, err := ts.playerRepository.GetPlayerOutOfTransferList(accountId, playerId)

	if err != nil {
		return 0, err
	}

	transferId, err = ts.transferRepository.NewTransfer(playerId, askedPrice, player.MarketValue)

	if err != nil {
		return 0, err
	}

	return transferId, err
}

func (ts *TransferService) ConfirmTransfer(accountId, transferId int) error {
	transfer, err := ts.transferRepository.GetTransfer(transferId)

	if err != nil {
		return err
	}

	buyer, err := ts.teamRepository.GetTeamByAccountId(accountId)

	if err != nil {
		return err
	}

	if *transfer.Player.TeamId == buyer.Id {
		return errors.New("destination team cannot be the same as the origin team")
	}

	if buyer.AvailableCash < transfer.AskedPrice {
		return errors.New("insufficient funds")
	}

	newMarketValue := (transfer.AskedPrice * (rand.Intn(191) + 110)) / 100
	return ts.transferRepository.ConfirmTransfer(transferId, buyer.Id, newMarketValue)
}

func (ts *TransferService) GetTransfers() ([]domain.Transfer, error) {
	transfers, err := ts.transferRepository.FindTransfers()
	return transfers, err
}

func (ts *TransferService) GetTransfer(transferId int) (*domain.Transfer, error) {
	transfer, err := ts.transferRepository.GetTransfer(transferId)
	return &transfer, err
}

func (ts *TransferService) UpdateTransfer(accountId int, transferId int, patchJSON []byte) (*domain.Transfer, error) {
	transfer, err := ts.GetTransfer(transferId)

	if err != nil {
		return nil, err
	}

	patch, err := jsonpatch.DecodePatch(patchJSON)

	if err != nil {
		return nil, err
	}

	transferJSON, _ := json.Marshal(transfer)
	patched, err := patch.Apply(transferJSON)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(patched, transfer); err != nil {
		return nil, err
	}

	err = ts.transferRepository.UpdateTransfer(accountId, transfer)

	if err != nil {
		return nil, err
	}

	return transfer, nil
}
