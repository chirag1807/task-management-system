package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/chirag1807/task-management-system/api/model/request"
	"github.com/chirag1807/task-management-system/api/model/response"
	"github.com/chirag1807/task-management-system/api/service"
	"github.com/chirag1807/task-management-system/constant"
	errorhandling "github.com/chirag1807/task-management-system/error"
	"github.com/chirag1807/task-management-system/utils"
	"github.com/go-chi/chi/v5"
)

type TeamController interface {
	CreateTeam(w http.ResponseWriter, r *http.Request)
	AddMembersToTeam(w http.ResponseWriter, r *http.Request)
	RemoveMembersFromTeam(w http.ResponseWriter, r *http.Request)
	GetAllTeams(w http.ResponseWriter, r *http.Request)
	GetTeamMembers(w http.ResponseWriter, r *http.Request)
	LeftTeam(w http.ResponseWriter, r *http.Request)
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
	var team request.CreateTeam

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errorhandling.SendErrorResponse(w, errorhandling.ReadBodyError)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &team)
	if err != nil {
		errorhandling.SendErrorResponse(w, errorhandling.ReadDataError)
		return
	}

	log.Println(team.TeamDetails, team.TeamMembers)
	err = t.teamService.CreateTeam(team.TeamDetails, team.TeamMembers)
	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}
	response := response.SuccessResponse{
		Message: constant.TEAM_CREATED,
	}
	utils.SendSuccessResponse(w, http.StatusOK, response)
}

func (t teamController) AddMembersToTeam(w http.ResponseWriter, r *http.Request) {
	var teamMembersToAdd request.TeamMembers

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errorhandling.SendErrorResponse(w, errorhandling.ReadBodyError)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &teamMembersToAdd)
	if err != nil {
		errorhandling.SendErrorResponse(w, errorhandling.ReadDataError)
		return
	}

	userId := r.Context().Value(constant.UserIdKey).(int64)
	err = t.teamService.AddMembersToTeam(userId, teamMembersToAdd)
	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}
	response := response.SuccessResponse{
		Message: constant.MEMBERS_ADDED_TO_TEAM,
	}
	utils.SendSuccessResponse(w, http.StatusOK, response)
}

func (t teamController) RemoveMembersFromTeam(w http.ResponseWriter, r *http.Request) {
	var teamMembersToRemove request.TeamMembers

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errorhandling.SendErrorResponse(w, errorhandling.ReadBodyError)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &teamMembersToRemove)
	if err != nil {
		errorhandling.SendErrorResponse(w, errorhandling.ReadDataError)
		return
	}

	userId := r.Context().Value(constant.UserIdKey).(int64)
	err = t.teamService.RemoveMembersFromTeam(userId, teamMembersToRemove)
	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}
	response := response.SuccessResponse{
		Message: constant.MEMBERS_REMOVED_FROM_TEAM,
	}
	utils.SendSuccessResponse(w, http.StatusOK, response)
}

func (t teamController) GetAllTeams(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(constant.UserIdKey).(int64)
	flag, err := strconv.Atoi(chi.URLParam(r, "Flag"))
	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}
	teams, err := t.teamService.GetAllTeams(userId, flag)
	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}
	utils.SendSuccessResponse(w, http.StatusOK, teams)
}

func (t teamController) GetTeamMembers(w http.ResponseWriter, r *http.Request) {
	teamID, err := strconv.ParseInt(chi.URLParam(r, "TeamID"), 10, 64)
	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}
	teams, err := t.teamService.GetTeamMembers(teamID)
	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}
	utils.SendSuccessResponse(w, http.StatusOK, teams)
}

func (t teamController) LeftTeam(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(constant.UserIdKey).(int64)
	teamID, err := strconv.ParseInt(chi.URLParam(r, "TeamID"), 10, 64)
	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}
	err = t.teamService.LeftTeam(userId, teamID)
	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}
	response := response.SuccessResponse{
		Message: constant.LEFT_TEAM,
	}
	utils.SendSuccessResponse(w, http.StatusOK, response)
}
