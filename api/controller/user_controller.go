package controller

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/chirag1807/task-management-system/api/model/request"
	"github.com/chirag1807/task-management-system/api/model/response"
	"github.com/chirag1807/task-management-system/api/service"
	"github.com/chirag1807/task-management-system/constant"
	errorhandling "github.com/chirag1807/task-management-system/error"
	"github.com/chirag1807/task-management-system/utils"
)

type UserController interface {
	GetAllPublicProfileUsers(w http.ResponseWriter, r *http.Request)
	GetMyDetails(w http.ResponseWriter, r *http.Request)
	UpdateUserProfile(w http.ResponseWriter, r *http.Request)
	SendOTPToUser(w http.ResponseWriter, r *http.Request)
	VerifyOTP(w http.ResponseWriter, r *http.Request)
	ResetUserPassword(w http.ResponseWriter, r *http.Request)
}

type userController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return userController{
		userService: userService,
	}
}

func (u userController) GetAllPublicProfileUsers(w http.ResponseWriter, r *http.Request) {
	publicProfileUsers, err := u.userService.GetAllPublicProfileUsers()
	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}
	response := response.Users{
		Users: publicProfileUsers,
	}
	utils.SendSuccessResponse(w, http.StatusOK, response)
}

func (u userController) GetMyDetails(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(constant.UserIdKey).(int64)
	userDetails, err := u.userService.GetMyDetails(userId)
	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}
	utils.SendSuccessResponse(w, http.StatusOK, userDetails)
}

func (u userController) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	var userToUpdate request.UpdateUser

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errorhandling.SendErrorResponse(w, errorhandling.ReadBodyError)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &userToUpdate)
	if err != nil {
		errorhandling.SendErrorResponse(w, errorhandling.ReadDataError)
		return
	}

	if userToUpdate.Password != nil {
		hashedPassword, err := utils.HashPassword(*userToUpdate.Password)
		if err != nil {
			errorhandling.SendErrorResponse(w, err)
			return
		}
		userToUpdate.Password = &hashedPassword
	}

	userId := r.Context().Value(constant.UserIdKey).(int64)
	err = u.userService.UpdateUserProfile(userId, userToUpdate)
	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}

	response := response.SuccessResponse{
		Message: constant.USER_PROFILE_UPDATED,
	}
	utils.SendSuccessResponse(w, http.StatusOK, response)
}

func (u userController) SendOTPToUser(w http.ResponseWriter, r *http.Request) {

}

func (u userController) VerifyOTP(w http.ResponseWriter, r *http.Request) {

}

func (u userController) ResetUserPassword(w http.ResponseWriter, r *http.Request) {

}
