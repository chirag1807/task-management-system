package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/chirag1807/task-management-system/api/model/request"
	"github.com/chirag1807/task-management-system/api/model/response"
	"github.com/chirag1807/task-management-system/api/service"
	"github.com/chirag1807/task-management-system/config"
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

// UserRegistration registers a new user in the task manager application.
// @Summary Register User
// @Description UserRegistration API is made for registering a new user in the task manager application.
// @Accept json
// @Produce json
// @Tags auth
// @Param firstName formData string true "First name of the user"
// @Param lastName formData string true "Last name of the user"
// @Param bio formData string true "Bio of the user"
// @Param email formData string true "Email of the user"
// @Param password formData string true "Password of the user"
// @Param profile formData string true "Profile of the user (Public, Private)"
// @Success 200 {object} response.SuccessResponse "User created successfully."
// @Failure 400 {object} errorhandling.CustomError "Bad request."
// @Failure 409 {object} errorhandling.CustomError "Duplicate email found."
// @Failure 500 {object} errorhandling.CustomError "Internal server error."
// @Router /api/auth/user-registration [post]
func (a authController) UserRegistration(w http.ResponseWriter, r *http.Request) {

	var requestParams = map[string]string{
		constant.FirstNameKey: "string|minLen:2|required",
		constant.LastNameKey:  "string|minLen:2|required",
		constant.BioKey:       "string|minLen:6|required",
		constant.EmailKey:     `string|regex:^[\w.%+-]+@[\w.-]+\.[a-zA-Z]{2,}$|required`,
		constant.PasswordKey:  "string|minLen:8|required",
		constant.ProfileKey:   "string|in:Public,Private|required",
	}
	var userRequest request.User
	err, invalidParamsMultiLineErrMsg := utils.ValidateParameters(r, &userRequest, &requestParams, nil, nil, nil, nil)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, err, utils.CreateErrorMessage())
		return
	}
	if invalidParamsMultiLineErrMsg != nil {
		errorhandling.SendErrorResponse(r, w, invalidParamsMultiLineErrMsg, "")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, errorhandling.ReadBodyError, "")
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &userRequest)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, errorhandling.ReadDataError, "")
		return
	}

	hashedPassword, err := utils.HashPassword(userRequest.Password)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, err, utils.CreateErrorMessage())
		return
	}
	userRequest.Password = hashedPassword

	userId, err := a.authService.UserRegistration(userRequest)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, err, utils.CreateErrorMessage())
		return
	}

	response := response.SuccessResponse{
		Message: constant.USER_REGISTRATION_SUCCEED,
		ID:      &userId,
	}
	config.LoggerInstance.Info(constant.USER_REGISTRATION_SUCCEED)
	utils.SendSuccessResponse(w, http.StatusOK, response)
}

// UserLogin login the user in task manager application.
// @Summary Login User
// @Description UserLogin API is made for login the user in task manager application.
// @Accept json
// @Produce json
// @Tags auth
// @Param email formData string true "Email of the user"
// @Param password formData string true "Password of the user"
// @Success 200 {object} response.UserWithTokens "User login done successfully."
// @Failure 400 {object} errorhandling.CustomError "Bad request."
// @Failure 401 {object} errorhandling.CustomError "Password not matched."
// @Failure 404 {object} errorhandling.CustomError "User not found."
// @Failure 500 {object} errorhandling.CustomError "Internal server error."
// @Router /api/auth/user-login [post]
func (a authController) UserLogin(w http.ResponseWriter, r *http.Request) {
	var requestParams = map[string]string{
		constant.EmailKey:    `string|regex:^[\w.%+-]+@[\w.-]+\.[a-zA-Z]{2,}$|required`,
		constant.PasswordKey: "string|minLen:8|required",
	}
	var userLoginRequest request.User

	err, invalidParamsMultiLineErrMsg := utils.ValidateParameters(r, &userLoginRequest, &requestParams, nil, nil, nil, nil)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, err, utils.CreateErrorMessage())
		return
	}
	if invalidParamsMultiLineErrMsg != nil {
		errorhandling.SendErrorResponse(r, w, invalidParamsMultiLineErrMsg, "")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, errorhandling.ReadBodyError, "")
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &userLoginRequest)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, errorhandling.ReadDataError, "")
		return
	}

	var user response.User
	var refreshToken string
	user, refreshToken, err = a.authService.UserLogin(userLoginRequest)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, err, utils.CreateErrorMessage())
		return
	}

	accessToken, err := utils.CreateJWTToken(time.Now().Add(time.Hour*5), user.ID)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, err, utils.CreateErrorMessage())
		return
	}

	response := response.UserWithTokens{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	config.LoggerInstance.Info(constant.USER_LOGIN_SUCCEED)
	utils.SendSuccessResponse(w, http.StatusOK, response)
}

// ResetToken reset the access token of user.
// @Summary Reset Access Token
// @Description ResetToken API is made for reset the user's access token.
// @Produce json
// @Tags auth
// @Param Authorization header string true "Refresh Token" default(Bearer <refresh_token>)
// @Success 200 {object} response.AccessToken "Token reset done successfully."
// @Failure 401 {object} errorhandling.CustomError "Either refresh token not found or token is expired."
// @Failure 500 {object} errorhandling.CustomError "Internal server error."
// @Router /api/auth/reset-token [post]
func (a authController) ResetToken(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value(constant.TokenKey).(string)

	userId, err := a.authService.ResetToken(token)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, err, utils.CreateErrorMessage())
		return
	}

	accessToken, err := utils.CreateJWTToken(time.Now().Add(time.Hour*5), userId)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, errorhandling.RefreshTokenError, utils.CreateErrorMessage())
		return
	}

	response := response.AccessToken{
		AccessToken: accessToken,
	}
	config.LoggerInstance.Info(constant.TOKEN_RESET_SUCCEED)
	utils.SendSuccessResponse(w, http.StatusOK, response)
}
