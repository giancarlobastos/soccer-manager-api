package service

import (
	"encoding/json"
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

func (ts *TeamService) GetTeam(teamId int) (*domain.Team, error) {
	team, err := ts.teamRepository.GetTeamById(teamId)

	if err != nil {
		return nil, err
	}

	players, err := ts.playerRepository.GetPlayersByTeamId(team.Id)

	if err != nil {
		return nil, err
	}

	team.Players = players

	return team, nil
}

func (ts *TeamService) UpdateTeam(teamId int, patchJSON []byte) (*domain.Team, error) {
	team, err := ts.GetTeam(teamId)

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

	err = ts.teamRepository.UpdateTeam(team)

	if err != nil {
		return nil, err
	}

	return team, nil
}
