package service

import (
	"encoding/json"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/giancarlobastos/soccer-manager-api/domain"
	"github.com/giancarlobastos/soccer-manager-api/repository"
)

type PlayerService struct {
	playerRepository *repository.PlayerRepository
}

func NewPlayerService(pr *repository.PlayerRepository) *PlayerService {
	return &PlayerService{
		playerRepository: pr,
	}
}

func (ps *PlayerService) GetPlayer(playerId int) (*domain.Player, error) {
	player, err := ps.playerRepository.GetPlayer(playerId)
	return &player, err
}

func (ps *PlayerService) UpdatePlayer(playerId int, patchJSON []byte) (*domain.Player, error) {
	player, err := ps.GetPlayer(playerId)

	if err != nil {
		return nil, err
	}

	patch, err := jsonpatch.DecodePatch(patchJSON)

	if err != nil {
		return nil, err
	}

	playerJSON, _ := json.Marshal(player)
	patched, err := patch.Apply(playerJSON)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(patched, player); err != nil {
		return nil, err
	}

	err = ps.playerRepository.UpdatePlayer(player)

	if err != nil {
		return nil, err
	}

	return player, nil
}
