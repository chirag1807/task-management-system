package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/chirag1807/task-management-system/api/model/request"
	"github.com/chirag1807/task-management-system/api/model/response"
	"github.com/chirag1807/task-management-system/api/service"
	"github.com/chirag1807/task-management-system/api/validation"
	"github.com/chirag1807/task-management-system/constant"
	errorhandling "github.com/chirag1807/task-management-system/error"
	"github.com/chirag1807/task-management-system/utils"
)

type AuthController interface {
	UserRegistration(w http.ResponseWriter, r *http.Request)
	UserLogin(w http.ResponseWriter, r *http.Request)
	ResetToken(w http.ResponseWriter, r *http.Request)
}

type authController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) AuthController {
	return authController{
		authService: authService,
	}
}

func (a authController) UserRegistration(w http.ResponseWriter, r *http.Request) {
	var requestParams = map[string]string{
		constant.FirstNameKey: "string|minLen:2",
		constant.LastNameKey:  "string|minLen:2",
		constant.BioKey:       "string|minLen:6",
		constant.EmailKey:     `string`,
		constant.PasswordKey:  "string|minLen:8",
		constant.ProfileKey:   "string|in:Public,Private",
	}
	var userRequest request.User
	err, invalidParamsMultiLineErrMsg, invalidParamsErrMsg := validation.ValidateParameters(r, &userRequest, &requestParams, nil, nil, nil, nil)

	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}
	log.Println(err, invalidParamsMultiLineErrMsg, invalidParamsErrMsg)

	if invalidParamsErrMsg != nil {
		errorhandling.SendErrorResponse(w, invalidParamsErrMsg)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errorhandling.SendErrorResponse(w, errorhandling.ReadBodyError)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &userRequest)
	if err != nil {
		errorhandling.SendErrorResponse(w, errorhandling.ReadDataError)
		return
	}

	isValidEmail := validation.EmailValidation(userRequest.Email)
	if !isValidEmail {
		errorhandling.SendErrorResponse(w, errorhandling.EmailvalidationError)
		return
	}

	hashedPassword, err := utils.HashPassword(userRequest.Password)
	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}
	userRequest.Password = hashedPassword

	userId, err := a.authService.UserRegistration(userRequest)
	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}

	response := response.SuccessResponse{
		Message: constant.USER_REGISTRATION_SUCCEED,
		ID:      &userId,
	}
	utils.SendSuccessResponse(w, http.StatusOK, response)
}

func (a authController) UserLogin(w http.ResponseWriter, r *http.Request) {
	var userLoginRequest request.User

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errorhandling.SendErrorResponse(w, errorhandling.ReadBodyError)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &userLoginRequest)
	if err != nil {
		errorhandling.SendErrorResponse(w, errorhandling.ReadDataError)
		return
	}

	isEmail := validation.EmailValidation(userLoginRequest.Email)
	if !isEmail {
		errorhandling.SendErrorResponse(w, errorhandling.EmailvalidationError)
		return
	}

	var user response.User
	var refreshToken string
	user, refreshToken, err = a.authService.UserLogin(userLoginRequest)

	if err != nil {
		errorhandling.SendErrorResponse(w, errorhandling.LoginFailedError)
		return
	}

	accessToken, err := utils.CreateJWTToken(time.Now().Add(time.Hour*5), user.ID)
	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}

	response := response.UserWithTokens{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	utils.SendSuccessResponse(w, http.StatusOK, response)
}

func (a authController) ResetToken(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value(constant.TokenKey).(string)

	userId, err := a.authService.ResetToken(token)
	if err != nil {
		errorhandling.SendErrorResponse(w, errorhandling.RefreshTokenError)
		return
	}

	accessToken, err := utils.CreateJWTToken(time.Now().Add(time.Hour*5), userId)
	if err != nil {
		errorhandling.SendErrorResponse(w, errorhandling.RefreshTokenError)
		return
	}

	response := response.AccessToken{
		AccessToken: accessToken,
	}
	utils.SendSuccessResponse(w, http.StatusOK, response)
}
