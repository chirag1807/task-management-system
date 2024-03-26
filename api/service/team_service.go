package service

import (
	"github.com/chirag1807/task-management-system/api/model/request"
	"github.com/chirag1807/task-management-system/api/model/response"
	"github.com/chirag1807/task-management-system/api/repository"
)

type TeamService interface {
	CreateTeam(teamToCreate request.Team, teamMembers request.TeamMembers) (int64, error)
	AddMembersToTeam(teamCreatedBy int64, teamMembersToAdd request.TeamMembers) error
	RemoveMembersFromTeam(teamCreatedBy int64, teamMembersToRemove request.TeamMembers) error
	GetAllTeams(userID int64, flag int, queryParams request.TeamQueryParams) ([]response.Team, error)
	GetTeamMembers(teamID int64, queryParams request.TeamQueryParams) ([]response.User, error)
	LeftTeam(userID int64, teamID int64) (error)
}

type teamService struct {
	teamRepository repository.TeamRepository
}

func NewTeamService(teamRepository repository.TeamRepository) TeamService {
	return teamService{
		teamRepository: teamRepository,
	}
}

func (t teamService) CreateTeam(teamToCreate request.Team, teamMembers request.TeamMembers) (int64, error) {
	return t.teamRepository.CreateTeam(teamToCreate, teamMembers)
}

func (t teamService) AddMembersToTeam(teamCreatedBy int64, teamMembersToAdd request.TeamMembers) error {
	return t.teamRepository.AddMembersToTeam(teamCreatedBy, teamMembersToAdd)
}

func (t teamService) RemoveMembersFromTeam(teamCreatedBy int64, teamMembersToRemove request.TeamMembers) error {
	return t.teamRepository.RemoveMembersFromTeam(teamCreatedBy, teamMembersToRemove)
}

func (t teamService) GetAllTeams(userID int64, flag int, queryParams request.TeamQueryParams) ([]response.Team, error) {
	return t.teamRepository.GetAllTeams(userID, flag, queryParams)
}

func (t teamService) GetTeamMembers(teamID int64, queryParams request.TeamQueryParams) ([]response.User, error) {
	return t.teamRepository.GetTeamMembers(teamID, queryParams)
}

func (t teamService) LeftTeam(userID int64, teamID int64) (error) {
	return t.teamRepository.LeftTeam(userID, teamID)
}
