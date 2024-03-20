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
	"github.com/chirag1807/task-management-system/api/validation"
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
		constant.TeamMembersIdKey: "slice|number|required",
	}
	var team request.CreateTeam

	err, invalidParamsMultiLineErrMsg, invalidParamsErrMsg := validation.ValidateParameters(r, &team, &requestParams, nil, nil, nil, nil)

	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}
	log.Println(err, invalidParamsMultiLineErrMsg, invalidParamsErrMsg)

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
	team.TeamMembers.MemberID = append(team.TeamMembers.MemberID, userId)
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
	var requestParams = map[string]string{
		constant.TeamIdKey:       "number|required",
		constant.TeamMemberIdKey: "slice|number|required",
	}
	var teamMembersToAdd request.TeamMembers

	err, invalidParamsMultiLineErrMsg, invalidParamsErrMsg := validation.ValidateParameters(r, &teamMembersToAdd, &requestParams, nil, nil, nil, nil)

	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}
	log.Println(err, invalidParamsMultiLineErrMsg, invalidParamsErrMsg)

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
	utils.SendSuccessResponse(w, http.StatusOK, response)
}

func (t teamController) RemoveMembersFromTeam(w http.ResponseWriter, r *http.Request) {
	var requestParams = map[string]string{
		constant.TeamIdKey:       "number|required",
		constant.TeamMemberIdKey: "slice|number|required",
	}
	var teamMembersToRemove request.TeamMembers

	err, invalidParamsMultiLineErrMsg, invalidParamsErrMsg := validation.ValidateParameters(r, &teamMembersToRemove, &requestParams, nil, nil, nil, nil)

	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}
	log.Println(err, invalidParamsMultiLineErrMsg, invalidParamsErrMsg)

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

	err, invalidParamsMultiLineErrMsg, invalidParamsErrMsg := validation.ValidateParameters(r, &teamQueryParams, nil, nil, &queryParams, &queryParamFilters, nil)
	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}
	log.Println(err, invalidParamsMultiLineErrMsg, invalidParamsErrMsg)

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

	err, invalidParamsMultiLineErrMsg, invalidParamsErrMsg := validation.ValidateParameters(r, &teamQueryParams, nil, nil, &queryParams, &queryParamFilters, nil)
	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}
	log.Println(err, invalidParamsMultiLineErrMsg, invalidParamsErrMsg)

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
	utils.SendSuccessResponse(w, http.StatusOK, response)
}
