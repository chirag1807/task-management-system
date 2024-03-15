package controller

import (
	"net/http"

	"github.com/chirag1807/task-management-system/api/service"
)

type UserController interface{
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

func (u userController) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {

}

func (u userController) SendOTPToUser(w http.ResponseWriter, r *http.Request) {
	
}

func (u userController) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	
}

func (u userController) ResetUserPassword(w http.ResponseWriter, r *http.Request) {
	
}