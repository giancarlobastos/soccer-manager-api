package service

import (
	"encoding/json"
	"errors"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/giancarlobastos/soccer-manager-api/domain"
	"github.com/giancarlobastos/soccer-manager-api/repository"
)

type TeamService struct {
	teamRepository   *repository.TeamRepository
	playerRepository *repository.PlayerRepository
}

func NewTeamService(tr *repository.TeamRepository, pr *repository.PlayerRepository) *TeamService {
	return &TeamService{
		teamRepository:   tr,
		playerRepository: pr,
	}
}

func (ts *TeamService) GetTeam(accountId, teamId int) (*domain.Team, error) {
	team, err := ts.teamRepository.GetTeamByAccountId(accountId)

	if err != nil {
		return nil, err
	}

	if team.Id != teamId {
		return nil, errors.New("team not owned by account id")
	}

	players, err := ts.playerRepository.GetPlayersByTeamId(team.Id)

	if err != nil {
		return nil, err
	}

	team.Players = players

	return team, nil
}

func (ts *TeamService) UpdateTeam(accountId, teamId int, patchJSON []byte) (*domain.Team, error) {
	team, err := ts.GetTeam(accountId, teamId)

	if err != nil {
		return nil, err
	}

	patch, err := jsonpatch.DecodePatch(patchJSON)

	if err != nil {
		return nil, err
	}

	teamJSON, _ := json.Marshal(team)
	patched, err := patch.Apply(teamJSON)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(patched, team); err != nil {
		return nil, err
	}

	err = ts.teamRepository.UpdateTeam(accountId, team)

	if err != nil {
		return nil, err
	}

	return team, nil
}
