package controller

import (
	"net/http"

	"github.com/chirag1807/task-management-system/api/service"
)

type TeamController interface{
	CreateTeam(w http.ResponseWriter, r *http.Request)
	AddMembersToTeam(w http.ResponseWriter, r *http.Request)
	RemoveMembersFromTeam(w http.ResponseWriter, r *http.Request)
	GetAllTeams(w http.ResponseWriter, r *http.Request)
	GetTeamMembers(w http.ResponseWriter, r *http.Request)
}

type teamController struct {
	teamService service.TeamService
}

func NewTeamController(teamService service.TeamService) TeamController {
	return teamController{
		teamService: teamService,
	}
}

func (t teamController) CreateTeam(w http.ResponseWriter, r *http.Request) {
	// t.teamService.CreateTeam(team, teamMembers)
}

func (t teamController) AddMembersToTeam(w http.ResponseWriter, r *http.Request) {
	// t.teamService.AddMembersToTeam(teamMembersToAdd)
}

func (t teamController) RemoveMembersFromTeam(w http.ResponseWriter, r *http.Request) {
	// t.teamService.RemoveMembersFromTeam(teamMembersToRemove)
}

func (t teamController) GetAllTeams(w http.ResponseWriter, r *http.Request) {
	// t.teamService.GetAllTeams(userID, flag)
}

func (t teamController) GetTeamMembers(w http.ResponseWriter, r *http.Request) {
	// t.teamService.GetTeamMembers(teamID)
}
