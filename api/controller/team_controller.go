package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/chirag1807/task-management-system/api/model/request"
	"github.com/chirag1807/task-management-system/api/model/response"
	"github.com/chirag1807/task-management-system/api/service"
	"github.com/chirag1807/task-management-system/config"
	"github.com/chirag1807/task-management-system/constant"
	errorhandling "github.com/chirag1807/task-management-system/error"
	"github.com/chirag1807/task-management-system/utils"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
)

type TeamController interface {
	CreateTeam(w http.ResponseWriter, r *http.Request)
	AddMembersToTeam(w http.ResponseWriter, r *http.Request)
	RemoveMembersFromTeam(w http.ResponseWriter, r *http.Request)
	GetAllTeams(w http.ResponseWriter, r *http.Request)
	GetTeamMembers(w http.ResponseWriter, r *http.Request)
	LeaveTeam(w http.ResponseWriter, r *http.Request)
}

type teamController struct {
	teamService service.TeamService
}

func NewTeamController(teamService service.TeamService) TeamController {
	return teamController{
		teamService: teamService,
	}
}

// CreateTeam creates a new team.
// @Summary Create New Team
// @Description CreateTeam API is made for creating a new team in the task manager application.
// @Accept json
// @Produce json
// @Tags teams
// @Param Authorization header string true "Access Token" default(Bearer <access_token>)
// @Param details formData request.Team true "Team name and profile"
// @Param members formData []int64 true "Ids of user who will be added to the team."
// @Success 200 {object} response.SuccessResponse "Team created successfully."
// @Failure 400 {object} errorhandling.CustomError "Bad request"
// @Failure 401 {object} errorhandling.CustomError "Either refresh token not found or token is expired."
// @Failure 500 {object} errorhandling.CustomError "Internal server error."
// @Router /api/v1/teams [post]
func (t teamController) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var team request.CreateTeam

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, errorhandling.ReadBodyError, constant.EMPTY_STRING)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &team)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, errorhandling.ReadDataError, constant.EMPTY_STRING)
		return
	}

	r.Body = io.NopCloser(bytes.NewReader(body))

	err = utils.Validate.Struct(team)
	if err != nil {
		errorhandling.HandleInvalidRequestData(w, r, err, utils.Translator)
		return
	}

	userId := r.Context().Value(constant.UserIdKey).(int64)
	team.Details.CreatedBy = userId
	if team.Details.Privacy == nil {
		defaultTeamPrivacy := "PUBLIC"
		team.Details.Privacy = &defaultTeamPrivacy
	}
	team.Members = append(team.Members, userId)

	teamId, err := t.teamService.CreateTeam(team.Details, team.Members)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, err, utils.CreateErrorMessage())
		return
	}
	response := response.SuccessResponse{
		Code:    http.StatusText(http.StatusOK),
		Message: constant.TEAM_CREATED,
		ID:      &teamId,
	}
	config.LoggerInstance.Info(constant.TEAM_CREATED)
	utils.SendSuccessResponse(w, http.StatusOK, response)
}

// AddMembersToTeam adds members to a team.
// @Summary Add members to a team
// @Description Add members to a team based on provided parameters
// @Accept json
// @Produce json
// @Tags teams
// @Param TeamID path int64 true "Team ID"
// @Param Authorization header string true "Access Token" default(Bearer <access_token>)
// @Param memberIds formData []int64 true "Array of member IDs to add to the team"
// @Success 200 {object} response.SuccessResponse "Members added successfully"
// @Failure 400 {object} errorhandling.CustomError "Bad request"
// @Failure 401 {object} errorhandling.CustomError "Either refresh token not found or token is expired."
// @Failure 403 {object} errorhandling.CustomError "Not allowed to add members."
// @Failure 409 {object} errorhandling.CustomError "Member already exist."
// @Failure 500 {object} errorhandling.CustomError "Internal server error."
// @Router /api/v1/teams/members [put]
func (t teamController) AddMembersToTeam(w http.ResponseWriter, r *http.Request) {
	var teamMembersToAdd request.TeamMembersWithTeamID

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, errorhandling.ReadBodyError, constant.EMPTY_STRING)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &teamMembersToAdd)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, errorhandling.ReadDataError, constant.EMPTY_STRING)
		return
	}

	teamId, err := strconv.ParseInt(chi.URLParam(r, constant.TEAM_ID), 10, 64)
	if err != nil {
		if strings.Contains(err.Error(), constant.URL_PARAM_CONVERT_ERROR) {
			errorhandling.SendErrorResponse(r, w, errorhandling.ProvideValidParams, constant.EMPTY_STRING)
			return
		}
		errorhandling.SendErrorResponse(r, w, err, utils.CreateErrorMessage())
		return
	}
	teamMembersToAdd.TeamID = teamId

	r.Body = io.NopCloser(bytes.NewReader(body))

	err = utils.Validate.Struct(teamMembersToAdd)
	if err != nil {
		errorhandling.HandleInvalidRequestData(w, r, err, utils.Translator)
		return
	}

	userId := r.Context().Value(constant.UserIdKey).(int64)
	err = t.teamService.AddMembersToTeam(userId, teamMembersToAdd)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, err, utils.CreateErrorMessage())
		return
	}
	response := response.SuccessResponse{
		Code:    http.StatusText(http.StatusOK),
		Message: constant.MEMBERS_ADDED_TO_TEAM,
	}
	config.LoggerInstance.Info(constant.MEMBERS_ADDED_TO_TEAM)
	utils.SendSuccessResponse(w, http.StatusOK, response)
}

// RemoveMembersFromTeam removes members from a team.
// @Summary Remove members from a team
// @Description Remove members from a team based on provided parameters
// @Accept json
// @Produce json
// @Tags teams
// @Param TeamID path int64 true "Team ID"
// @Param Authorization header string true "Access Token" default(Bearer <access_token>)
// @Param memberIds formData []int64 true "Array of member IDs to add to the team"
// @Success 200 {object} response.SuccessResponse "Members Removed successfully"
// @Failure 400 {object} errorhandling.CustomError "Bad request"
// @Failure 401 {object} errorhandling.CustomError "Either refresh token not found or token is expired."
// @Failure 403 {object} errorhandling.CustomError "Not allowed to add members."
// @Failure 500 {object} errorhandling.CustomError "Internal server error."
// @Router /api/v1/teams/members [delete]
func (t teamController) RemoveMembersFromTeam(w http.ResponseWriter, r *http.Request) {
	var teamMembersToRemove request.TeamMembersWithTeamID

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, errorhandling.ReadBodyError, constant.EMPTY_STRING)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &teamMembersToRemove)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, errorhandling.ReadDataError, constant.EMPTY_STRING)
		return
	}

	teamId, err := strconv.ParseInt(chi.URLParam(r, constant.TEAM_ID), 10, 64)
	if err != nil {
		if strings.Contains(err.Error(), constant.URL_PARAM_CONVERT_ERROR) {
			errorhandling.SendErrorResponse(r, w, errorhandling.ProvideValidParams, constant.EMPTY_STRING)
			return
		}
		errorhandling.SendErrorResponse(r, w, err, utils.CreateErrorMessage())
		return
	}
	teamMembersToRemove.TeamID = teamId

	r.Body = io.NopCloser(bytes.NewReader(body))

	err = utils.Validate.Struct(teamMembersToRemove)
	if err != nil {
		errorhandling.HandleInvalidRequestData(w, r, err, utils.Translator)
		return
	}

	userId := r.Context().Value(constant.UserIdKey).(int64)
	err = t.teamService.RemoveMembersFromTeam(userId, teamMembersToRemove)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, err, utils.CreateErrorMessage())
		return
	}
	response := response.SuccessResponse{
		Code:    http.StatusText(http.StatusOK),
		Message: constant.MEMBERS_REMOVED_FROM_TEAM,
	}
	config.LoggerInstance.Info(constant.MEMBERS_REMOVED_FROM_TEAM)
	utils.SendSuccessResponse(w, http.StatusOK, response)
}

// GetAllTeams fetches all teams of user.
// @Summary Get all teams
// @Description Get all teams of user based on query parameters
// @Produce json
// @Tags teams
// @Param Authorization header string true "Access Token" default(Bearer <access_token>)
// @Param createdByMe query bool true "return teams created by you if createdByMe set to true otherwise false."
// @Param limit query int false "Number of tasks to return per page (default 10)"
// @Param offset query int false "Offset for pagination (default 0)"
// @Param search query string false "Search term to filter tasks"
// @Param sortByCreatedAt query bool false "Sort tasks by create time (true for ascending, false for descending)"
// @Success 200 {object} []response.Team "Teams fetched successfully."
// @Failure 400 {object} errorhandling.CustomError "Bad request"
// @Failure 401 {object} errorhandling.CustomError "Either refresh token not found or token is expired."
// @Failure 422 {object} errorhandling.CustomError "Provide valid flag"
// @Failure 500 {object} errorhandling.CustomError "Internal server error"
// @Router /api/v1/teams/{Flag} [get]
func (t teamController) GetAllTeams(w http.ResponseWriter, r *http.Request) {
	var teamQueryParams request.TeamQueryParams

	decoder := schema.NewDecoder()
	err := decoder.Decode(&teamQueryParams, r.URL.Query())
	fmt.Println(err)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, errorhandling.ReadQueryParamsError, constant.EMPTY_STRING)
		return
	}

	err = utils.Validate.Struct(teamQueryParams)
	if err != nil {
		errorhandling.HandleInvalidRequestData(w, r, err, utils.Translator)
		return
	}

	if teamQueryParams.Limit == 0 {
		teamQueryParams.Limit = 10
	}

	userId := r.Context().Value(constant.UserIdKey).(int64)
	teams, err := t.teamService.GetAllTeams(userId, teamQueryParams)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, err, utils.CreateErrorMessage())
		return
	}
	utils.SendSuccessResponse(w, http.StatusOK, teams)
}

// GetTeamMembers fetches all members of the team.
// @Summary Get all team members
// @Description Get all members of team based on query parameters
// @Produce json
// @Tags teams
// @Param Authorization header string true "Access Token" default(Bearer <access_token>)
// @Param TeamID path int64 true "ID of team whose members you want."
// @Param limit query int false "Number of tasks to return per page (default 10)"
// @Param offset query int false "Offset for pagination (default 0)"
// @Success 200 {object} []response.User "Team members fetched successfully."
// @Failure 400 {object} errorhandling.CustomError "Bad request"
// @Failure 401 {object} errorhandling.CustomError "Either refresh token not found or token is expired."
// @Failure 500 {object} errorhandling.CustomError "Internal server error"
// @Router /api/v1/teams/{TeamID}/members [get]
func (t teamController) GetTeamMembers(w http.ResponseWriter, r *http.Request) {
	var teamQueryParams request.TeamQueryParams

	decoder := schema.NewDecoder()
	err := decoder.Decode(&teamQueryParams, r.URL.Query())
	if err != nil {
		errorhandling.SendErrorResponse(r, w, errorhandling.ReadQueryParamsError, constant.EMPTY_STRING)
		return
	}

	err = utils.Validate.Struct(teamQueryParams)
	if err != nil {
		errorhandling.HandleInvalidRequestData(w, r, err, utils.Translator)
		return
	}

	if teamQueryParams.Limit == 0 {
		teamQueryParams.Limit = 10
	}

	teamId, err := strconv.ParseInt(chi.URLParam(r, constant.TEAM_ID), 10, 64)
	if err != nil {
		if strings.Contains(err.Error(), constant.URL_PARAM_CONVERT_ERROR) {
			errorhandling.SendErrorResponse(r, w, errorhandling.ProvideValidParams, constant.EMPTY_STRING)
			return
		}
		errorhandling.SendErrorResponse(r, w, err, utils.CreateErrorMessage())
		return
	}
	teamMembers, err := t.teamService.GetTeamMembers(teamId, teamQueryParams)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, err, utils.CreateErrorMessage())
		return
	}
	utils.SendSuccessResponse(w, http.StatusOK, teamMembers)
}

// LeaveTeam removes user from particular team.
// @Summary Leave Team
// @Description Removes user from particular team
// @Produce json
// @Tags teams
// @Param Authorization header string true "Access Token" default(Bearer <access_token>)
// @Param TeamID path int64 true "ID of team whose members you want."
// @Success 200 {object} response.SuccessResponse "Team left successfully."
// @Failure 401 {object} errorhandling.CustomError "Either refresh token not found or token is expired or you are not a member of that team."
// @Failure 500 {object} errorhandling.CustomError "Internal server error"
// @Router /api/v1/teams/leave/{TeamID} [delete]
func (t teamController) LeaveTeam(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(constant.UserIdKey).(int64)
	teamId, err := strconv.ParseInt(chi.URLParam(r, constant.TEAM_ID), 10, 64)
	if err != nil {
		if strings.Contains(err.Error(), constant.URL_PARAM_CONVERT_ERROR) {
			errorhandling.SendErrorResponse(r, w, errorhandling.ProvideValidParams, constant.EMPTY_STRING)
			return
		}
		errorhandling.SendErrorResponse(r, w, err, utils.CreateErrorMessage())
		return
	}
	err = t.teamService.LeaveTeam(userId, teamId)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, err, utils.CreateErrorMessage())
		return
	}
	response := response.SuccessResponse{
		Code:    http.StatusText(http.StatusOK),
		Message: constant.LEAVE_TEAM,
	}
	config.LoggerInstance.Info(constant.LEAVE_TEAM)
	utils.SendSuccessResponse(w, http.StatusOK, response)
}
