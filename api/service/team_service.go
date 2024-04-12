package service

import (
	"github.com/chirag1807/task-management-system/api/model/request"
	"github.com/chirag1807/task-management-system/api/model/response"
	"github.com/chirag1807/task-management-system/api/repository"
)

type TeamService interface {
	CreateTeam(teamToCreate request.Team, teamMembers []int64) (int64, error)
	AddMembersToTeam(teamCreatedBy int64, teamMembersToAdd request.TeamMembersWithTeamID) error
	RemoveMembersFromTeam(teamCreatedBy int64, teamMembersToRemove request.TeamMembersWithTeamID) error
	GetAllTeams(userID int64, queryParams request.TeamQueryParams) ([]response.Team, error)
	GetTeamMembers(teamId int64, queryParams request.TeamQueryParams) ([]response.User, error)
	LeaveTeam(userID int64, teamId int64) (error)
}

type teamService struct {
	teamRepository repository.TeamRepository
}

func NewTeamService(teamRepository repository.TeamRepository) TeamService {
	return teamService{
		teamRepository: teamRepository,
	}
}

func (t teamService) CreateTeam(teamToCreate request.Team, teamMembers []int64) (int64, error) {
	return t.teamRepository.CreateTeam(teamToCreate, teamMembers)
}

func (t teamService) AddMembersToTeam(teamCreatedBy int64, teamMembersToAdd request.TeamMembersWithTeamID) error {
	return t.teamRepository.AddMembersToTeam(teamCreatedBy, teamMembersToAdd)
}

func (t teamService) RemoveMembersFromTeam(teamCreatedBy int64, teamMembersToRemove request.TeamMembersWithTeamID) error {
	return t.teamRepository.RemoveMembersFromTeam(teamCreatedBy, teamMembersToRemove)
}

func (t teamService) GetAllTeams(userID int64, queryParams request.TeamQueryParams) ([]response.Team, error) {
	return t.teamRepository.GetAllTeams(userID, queryParams)
}

func (t teamService) GetTeamMembers(teamId int64, queryParams request.TeamQueryParams) ([]response.User, error) {
	return t.teamRepository.GetTeamMembers(teamId, queryParams)
}

func (t teamService) LeaveTeam(userID int64, teamId int64) (error) {
	return t.teamRepository.LeaveTeam(userID, teamId)
}
