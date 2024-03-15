package service

import (
	"github.com/chirag1807/task-management-system/api/model/request"
	"github.com/chirag1807/task-management-system/api/model/response"
	"github.com/chirag1807/task-management-system/api/repository"
)

type TeamService interface {
	CreateTeam(teamToCreate request.Team, teamMembers request.TeamMembers) error
	AddMembersToTeam(teamMembersToAdd request.TeamMembers) error
	RemoveMembersFromTeam(teamMembersToRemove request.TeamMembers) error
	GetAllTeams(userID int64, flag string) ([]response.Team, error)
	GetTeamMembers(teamID int64) (response.TeamMembers, error)
}

type teamService struct {
	teamRepository repository.TeamRepository
}

func NewTeamService(teamRepository repository.TeamRepository) TeamService {
	return teamService{
		teamRepository: teamRepository,
	}
}

func (t teamService) CreateTeam(teamToCreate request.Team, teamMembers request.TeamMembers) error {
	return t.teamRepository.CreateTeam(teamToCreate, teamMembers)
}

func (t teamService) AddMembersToTeam(teamMembersToAdd request.TeamMembers) error {
	return t.teamRepository.AddMembersToTeam(teamMembersToAdd)
}

func (t teamService) RemoveMembersFromTeam(teamMembersToRemove request.TeamMembers) error {
	return t.teamRepository.RemoveMembersFromTeam(teamMembersToRemove)
}

func (t teamService) GetAllTeams(userID int64, flag string) ([]response.Team, error) {
	return t.teamRepository.GetAllTeams(userID, flag)
}

func (t teamService) GetTeamMembers(teamID int64) (response.TeamMembers, error) {
	return t.teamRepository.GetTeamMembers(teamID)
}
