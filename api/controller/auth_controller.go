package controller

import (
	"bytes"
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

func init() {
	utils.InitReqDataValidationTranslation()
}

type AuthController interface {
	UserRegistration(w http.ResponseWriter, r *http.Request)
	UserLogin(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)
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
// @Param privacy formData string true "privacy of the user (PUBLIC, PRIVATE)"
// @Success 200 {object} response.SuccessResponse "User created successfully."
// @Failure 400 {object} errorhandling.CustomError "Bad request."
// @Failure 409 {object} errorhandling.CustomError "Duplicate email found."
// @Failure 500 {object} errorhandling.CustomError "Internal server error."
// @Router /api/v1/auth/registration [post]
func (a authController) UserRegistration(w http.ResponseWriter, r *http.Request) {
	var userRequest request.User

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, errorhandling.ReadBodyError, constant.EMPTY_STRING)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &userRequest)
	if err != nil {
		errorhandling.HandleJSONUnmarshlError(r, w, err)
		return
	}

	r.Body = io.NopCloser(bytes.NewReader(body))

	err = utils.Validate.Struct(userRequest)
	if err != nil {
		errorhandling.HandleInvalidRequestData(w, r, err, utils.Translator)
		return
	}

	if userRequest.Password != userRequest.ConfirmPassword {
		errorhandling.SendErrorResponse(r, w, errorhandling.PasswordConfirmPasswordNotMatched, constant.EMPTY_STRING)
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
		Code:    http.StatusText(http.StatusOK),
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
// @Router /api/v1/auth/login [post]
func (a authController) UserLogin(w http.ResponseWriter, r *http.Request) {
	var userLoginRequest request.UserCredentials

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, errorhandling.ReadBodyError, constant.EMPTY_STRING)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &userLoginRequest)
	if err != nil {
		errorhandling.HandleJSONUnmarshlError(r, w, err)
		return
	}

	r.Body = io.NopCloser(bytes.NewReader(body))

	err = utils.Validate.Struct(userLoginRequest)
	if err != nil {
		errorhandling.HandleInvalidRequestData(w, r, err, utils.Translator)
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
		Code:         http.StatusText(http.StatusOK),
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
// @Success 200 {object} response.Tokens "Token refresh done successfully."
// @Failure 401 {object} errorhandling.CustomError "Either refresh token not found or token is expired."
// @Failure 500 {object} errorhandling.CustomError "Internal server error."
// @Router /api/v1/auth/reset-token [post]
func (a authController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value(constant.TokenKey).(string)

	userId, refreshToken, err := a.authService.RefreshToken(token)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, err, utils.CreateErrorMessage())
		return
	}

	accessToken, err := utils.CreateJWTToken(time.Now().Add(time.Hour*5), userId)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, errorhandling.RefreshTokenError, utils.CreateErrorMessage())
		return
	}

	response := response.Tokens{
		Code:         http.StatusText(http.StatusOK),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	config.LoggerInstance.Info(constant.TOKEN_RESET_SUCCEED)
	utils.SendSuccessResponse(w, http.StatusOK, response)
}
