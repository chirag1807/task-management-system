package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

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
	var requestParams = map[string]string{
		constant.TeamNameKey:      "string|minLen:3|maxLen:15|required",
		constant.TeamProfileKey:   "string|in:Public,Private",
		constant.TeamMembersKey:   "required",
		constant.TeamMembersIdKey: "slice",
	}
	var team request.CreateTeam

	err, invalidParamsMultiLineErrMsg := utils.ValidateParameters(r, &team, &requestParams, nil, nil, nil, nil)

	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}

	if invalidParamsMultiLineErrMsg != nil {
		errorhandling.SendErrorResponse(w, invalidParamsMultiLineErrMsg)
		return
	}

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

	userId := r.Context().Value(constant.UserIdKey).(int64)
	team.TeamDetails.CreatedBy = userId
	if team.TeamDetails.TeamProfile == nil {
		defaultTeamProfile := "Public"
		team.TeamDetails.TeamProfile = &defaultTeamProfile
	}
	team.TeamMembers.MemberID = append(team.TeamMembers.MemberID, userId)

	teamId, err := t.teamService.CreateTeam(team.TeamDetails, team.TeamMembers)
	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}
	response := response.SuccessResponse{
		Message: constant.TEAM_CREATED,
		ID:      &teamId,
	}
	log.Println("Team Created Successfully.")
	utils.SendSuccessResponse(w, http.StatusOK, response)
}

// AddMembersToTeam adds members to a team.
// @Summary Add members to a team
// @Description Add members to a team based on provided parameters
// @Accept json
// @Produce json
// @Tags teams
// @Param Authorization header string true "Access Token" default(Bearer <access_token>)
// @Param teamID formData int true "Team ID"
// @Param memberID formData array true "Array of member IDs to add to the team"
// @Success 200 {object} response.SuccessResponse "Members added successfully"
// @Failure 400 {object} errorhandling.CustomError "Bad request"
// @Failure 401 {object} errorhandling.CustomError "Either refresh token not found or token is expired."
// @Failure 403 {object} errorhandling.CustomError "Not allowed to add members."
// @Failure 409 {object} errorhandling.CustomError "Member already exist."
// @Failure 500 {object} errorhandling.CustomError "Internal server error."
// @Router /api/teams/add-members-to-team [put]
func (t teamController) AddMembersToTeam(w http.ResponseWriter, r *http.Request) {
	var requestParams = map[string]string{
		constant.TeamIdKey:       "number|required",
		constant.TeamMemberIdKey: "slice|required",
	}
	var teamMembersToAdd request.TeamMembers

	err, invalidParamsMultiLineErrMsg := utils.ValidateParameters(r, &teamMembersToAdd, &requestParams, nil, nil, nil, nil)

	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}

	if invalidParamsMultiLineErrMsg != nil {
		errorhandling.SendErrorResponse(w, invalidParamsMultiLineErrMsg)
		return
	}

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
	log.Println("Members Added to Team Successfully.")
	utils.SendSuccessResponse(w, http.StatusOK, response)
}

// RemoveMembersFromTeam removes members from a team.
// @Summary Remove members from a team
// @Description Remove members from a team based on provided parameters
// @Accept json
// @Produce json
// @Tags teams
// @Param Authorization header string true "Access Token" default(Bearer <access_token>)
// @Param teamID formData int true "Team ID"
// @Param memberID formData array true "Array of member IDs to add to the team"
// @Success 200 {object} response.SuccessResponse "Members Removed successfully"
// @Failure 400 {object} errorhandling.CustomError "Bad request"
// @Failure 401 {object} errorhandling.CustomError "Either refresh token not found or token is expired."
// @Failure 403 {object} errorhandling.CustomError "Not allowed to add members."
// @Failure 500 {object} errorhandling.CustomError "Internal server error."
// @Router /api/teams/remove-members-from-team [put]
func (t teamController) RemoveMembersFromTeam(w http.ResponseWriter, r *http.Request) {
	var requestParams = map[string]string{
		constant.TeamIdKey:       "number|required",
		constant.TeamMemberIdKey: "slice|required",
	}
	var teamMembersToRemove request.TeamMembers

	err, invalidParamsMultiLineErrMsg := utils.ValidateParameters(r, &teamMembersToRemove, &requestParams, nil, nil, nil, nil)

	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}

	if invalidParamsMultiLineErrMsg != nil {
		errorhandling.SendErrorResponse(w, invalidParamsMultiLineErrMsg)
		return
	}

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
	log.Println("Members Removed from Team Successfully.")
	utils.SendSuccessResponse(w, http.StatusOK, response)
}

func (t teamController) GetAllTeams(w http.ResponseWriter, r *http.Request) {
	var queryParams = map[string]string{
		constant.LimitKey:          "number|default:10",
		constant.OffsetKey:         "number|default:0",
		constant.SearchKey:         "string",
		constant.SortByCreateAtKey: "bool",
	}
	var queryParamFilters = map[string]string{
		constant.LimitKey:          "int",
		constant.OffsetKey:         "int",
		constant.SortByCreateAtKey: "bool",
	}

	var teamQueryParams request.TeamQueryParams

	err, invalidParamsMultiLineErrMsg := utils.ValidateParameters(r, &teamQueryParams, nil, nil, &queryParams, &queryParamFilters, nil)
	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}

	if invalidParamsMultiLineErrMsg != nil {
		errorhandling.SendErrorResponse(w, invalidParamsMultiLineErrMsg)
		return
	}

	userId := r.Context().Value(constant.UserIdKey).(int64)
	flag, err := strconv.Atoi(chi.URLParam(r, "Flag"))
	if err != nil {
		if strings.Contains(err.Error(), "strconv.Atoi: parsing") {
			errorhandling.SendErrorResponse(w, errorhandling.ProvideValidFlag)
			return
		}
		errorhandling.SendErrorResponse(w, err)
		return
	}
	if flag == 0 || flag == 1 {
		teams, err := t.teamService.GetAllTeams(userId, flag, teamQueryParams)
		if err != nil {
			errorhandling.SendErrorResponse(w, err)
			return
		}
		response := response.Teams{
			Teams: teams,
		}
		utils.SendSuccessResponse(w, http.StatusOK, response)
	} else {
		errorhandling.SendErrorResponse(w, errorhandling.ProvideValidFlag)
	}
}

func (t teamController) GetTeamMembers(w http.ResponseWriter, r *http.Request) {
	var queryParams = map[string]string{
		constant.LimitKey:  "number|default:10",
		constant.OffsetKey: "number|default:0",
	}
	var queryParamFilters = map[string]string{
		constant.LimitKey:  "int",
		constant.OffsetKey: "int",
	}

	var teamQueryParams request.TeamQueryParams

	err, invalidParamsMultiLineErrMsg := utils.ValidateParameters(r, &teamQueryParams, nil, nil, &queryParams, &queryParamFilters, nil)
	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}

	if invalidParamsMultiLineErrMsg != nil {
		errorhandling.SendErrorResponse(w, invalidParamsMultiLineErrMsg)
		return
	}

	teamID, err := strconv.ParseInt(chi.URLParam(r, "TeamID"), 10, 64)
	if err != nil {
		if strings.Contains(err.Error(), "strconv.ParseInt: parsing") {
			errorhandling.SendErrorResponse(w, errorhandling.ProvideValidParams)
			return
		}
		errorhandling.SendErrorResponse(w, err)
		return
	}
	teamMembers, err := t.teamService.GetTeamMembers(teamID, teamQueryParams)
	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}
	response := response.TeamMemberDetails{
		TeamMembers: teamMembers,
	}
	utils.SendSuccessResponse(w, http.StatusOK, response)
}

func (t teamController) LeftTeam(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(constant.UserIdKey).(int64)
	teamID, err := strconv.ParseInt(chi.URLParam(r, "TeamID"), 10, 64)
	if err != nil {
		if strings.Contains(err.Error(), "strconv.ParseInt: parsing") {
			errorhandling.SendErrorResponse(w, errorhandling.ProvideValidParams)
			return
		}
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
	log.Println("Team Left Successfully.")
	utils.SendSuccessResponse(w, http.StatusOK, response)
}
