package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/chirag1807/task-management-system/api/model/dto"
	"github.com/chirag1807/task-management-system/api/model/request"
	"github.com/chirag1807/task-management-system/api/model/response"
	"github.com/chirag1807/task-management-system/api/service"
	"github.com/chirag1807/task-management-system/config"
	"github.com/chirag1807/task-management-system/constant"
	errorhandling "github.com/chirag1807/task-management-system/error"
	"github.com/chirag1807/task-management-system/utils"
	"github.com/gorilla/schema"
)

type UserController interface {
	GetAllPublicPrivacyUsers(w http.ResponseWriter, r *http.Request)
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

// GetAllPublicPrivacyUsers fetches all public privacy users.
// @Summary Get all public privacy users
// @Description Get all public privacy users based on query parameters
// @Produce json
// @Tags users
// @Param Authorization header string true "Access Token" default(Bearer <access_token>)
// @Param Limit query int false "Number of users to return per page (default 10)"
// @Param Offset query int false "Offset for pagination (default 0)"
// @Param Search query string false "Search term to filter users"
// @Success 200 {object} []response.User "Public privacy users fetched successfully"
// @Failure 400 {object} errorhandling.CustomError "Bad request"
// @Failure 401 {object} errorhandling.CustomError "Either refresh token not found or token is expired."
// @Failure 500 {object} errorhandling.CustomError "Internal server error."
// @Router /api/v1/users/public-privacy [get]
func (u userController) GetAllPublicPrivacyUsers(w http.ResponseWriter, r *http.Request) {
	var userQueryParams request.UserQueryParams

	decoder := schema.NewDecoder()
	err := decoder.Decode(&userQueryParams, r.URL.Query())
	if err != nil {
		errorhandling.HandleSchemaDecodeError(r, w, err)
		return
	}

	err = utils.Validate.Struct(userQueryParams)
	if err != nil {
		errorhandling.HandleInvalidRequestData(w, r, err, utils.Translator)
		return
	}

	if userQueryParams.Limit == 0 {
		userQueryParams.Limit = 10
	}

	publicPrivacyUsers, err := u.userService.GetAllPublicPrivacyUsers(userQueryParams)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, err, utils.CreateErrorMessage())
		return
	}
	utils.SendSuccessResponse(w, http.StatusOK, publicPrivacyUsers)
}

// GetMyDetails fetches details of the authenticated user.
// @Summary Get details of the authenticated user
// @Description Get details of the authenticated user based on the authenticated user ID provided via token.
// @Produce json
// @Tags users
// @Param Authorization header string true "Access Token" default(Bearer <access_token>)
// @Success 200 {object} response.User "Successful response"
// @Failure 400 {object} errorhandling.CustomError "Bad request"
// @Failure 401 {object} errorhandling.CustomError "Either refresh token not found or token is expired."
// @Failure 500 {object} errorhandling.CustomError "Internal server error."
// @Router /api/v1/users/profile [get]
func (u userController) GetMyDetails(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(constant.UserIdKey).(int64)
	userDetails, err := u.userService.GetMyDetails(userId)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, err, utils.CreateErrorMessage())
		return
	}
	utils.SendSuccessResponse(w, http.StatusOK, userDetails)
}

// UpdateUserProfile updates a user's profile.
// @Summary Update User Profile
// @Description UpdateUserProfile API is made for updating a user's profile.
// @Accept json
// @Produce json
// @Tags users
// @Param Authorization header string true "Access Token" default(Bearer <access_token>)
// @Param firstName formData string false "First name of the user"
// @Param lastName formData string false "Last name of the user"
// @Param bio formData string false "Bio of the user"
// @Param email formData string false "Email of the user"
// @Param password formData string false "Password of the user"
// @Param privacy formData string false "Privacy of the user (PUBLIC, PRIVATE)"
// @Success 200 {object} response.SuccessResponse "User Updated successfully."
// @Failure 400 {object} errorhandling.CustomError "Bad request."
// @Failure 401 {object} errorhandling.CustomError "Either password not matched or token expired."
// @Failure 404 {object} errorhandling.CustomError "No user found."
// @Failure 409 {object} errorhandling.CustomError "Duplicate email found."
// @Failure 500 {object} errorhandling.CustomError "Internal server error."
// @Router /api/v1/users/profile [put]
func (u userController) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	var userToUpdate request.UpdateUser

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, errorhandling.ReadBodyError, constant.EMPTY_STRING)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &userToUpdate)
	if err != nil {
		errorhandling.HandleJSONUnmarshlError(r, w, err)
		return
	}

	r.Body = io.NopCloser(bytes.NewReader(body))

	err = utils.Validate.Struct(userToUpdate)
	if err != nil {
		errorhandling.HandleInvalidRequestData(w, r, err, utils.Translator)
		return
	}

	userId := r.Context().Value(constant.UserIdKey).(int64)
	if userToUpdate.Password != constant.EMPTY_STRING {
		err := u.userService.VerifyUserPassword(userToUpdate.Password, userId)
		if err != nil {
			errorhandling.SendErrorResponse(r, w, err, utils.CreateErrorMessage())
			return
		}
		hashedPassword, err := utils.HashPassword(userToUpdate.NewPassword)
		if err != nil {
			errorhandling.SendErrorResponse(r, w, err, utils.CreateErrorMessage())
			return
		}
		userToUpdate.Password = hashedPassword
	}

	err = u.userService.UpdateUserProfile(userId, userToUpdate)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, err, utils.CreateErrorMessage())
		return
	}

	response := response.SuccessResponse{
		Code: http.StatusText(http.StatusOK),
		Message: constant.USER_PROFILE_UPDATED,
	}
	config.LoggerInstance.Info(constant.USER_PROFILE_UPDATED)
	utils.SendSuccessResponse(w, http.StatusOK, response)
}

// SendOTPToUser sends an otp to user's email address.
// @Summary Sends an OTP
// @Description SendOTPToUser API is made for sending an otp to user's email address.
// @Accept json
// @Produce json
// @Tags users
// @Param email formData string true "Email of the user"
// @Success 200 {object} response.SuccessResponse "OTP sent successfully."
// @Failure 400 {object} errorhandling.CustomError "Bad request."
// @Failure 404 {object} errorhandling.CustomError "No Email found."
// @Failure 500 {object} errorhandling.CustomError "Internal server error."
// @Router /api/v1/users/send-otp [post]
func (u userController) SendOTPToUser(w http.ResponseWriter, r *http.Request) {
	var userEmail request.UserEmail

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, errorhandling.ReadBodyError, constant.EMPTY_STRING)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &userEmail)
	if err != nil {
		errorhandling.HandleJSONUnmarshlError(r, w, err)
		return
	}

	r.Body = io.NopCloser(bytes.NewReader(body))

	err = utils.Validate.Struct(userEmail)
	if err != nil {
		errorhandling.HandleInvalidRequestData(w, r, err, utils.Translator)
		return
	}

	min, max := 1001, 9999
	OTP := rand.Intn(max-min) + min
	fmt.Println(OTP)

	emailBody := utils.PrepareEmailBody(OTP)
	email := dto.Email{
		To:      userEmail.Email,
		Subject: "OTP Verification",
		Body:    emailBody,
	}

	otpExpireTime := time.Now().UTC().Add(time.Minute * 5)
	id, err := u.userService.SendOTPToUser(email, OTP, otpExpireTime)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, err, utils.CreateErrorMessage())
		return
	}

	config.LoggerInstance.Info(string(rune(OTP)))
	response := response.SuccessResponse{
		Code: http.StatusText(http.StatusOK),
		Message: constant.OTP_SENT,
		ID:      &id,
	}
	utils.SendSuccessResponse(w, http.StatusOK, response)
}

// VerifyOTP verifies an otp from user.
// @Summary Verifies an OTP
// @Description VerifyOTP API is made for verifying an otp from user.
// @Accept json
// @Produce json
// @Tags users
// @Param id formData int64 true "ID which you've received in response of SendOTPToUser API"
// @Param otp formData int true "OTP which user has entered"
// @Success 200 {object} response.SuccessResponse "OTP Verifies successfully."
// @Failure 400 {object} errorhandling.CustomError "Bad request."
// @Failure 401 {object} errorhandling.CustomError "OTP not matched."
// @Failure 410 {object} errorhandling.CustomError "OTP verification time expired."
// @Failure 500 {object} errorhandling.CustomError "Internal server error."
// @Router /api/v1/users/verify-otp [post]
func (u userController) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var otp request.OTP

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, errorhandling.ReadBodyError, constant.EMPTY_STRING)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &otp)
	if err != nil {
		errorhandling.HandleJSONUnmarshlError(r, w, err)
		return
	}

	r.Body = io.NopCloser(bytes.NewReader(body))

	err = utils.Validate.Struct(otp)
	if err != nil {
		errorhandling.HandleInvalidRequestData(w, r, err, utils.Translator)
		return
	}

	err = u.userService.VerifyOTP(otp)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, err, utils.CreateErrorMessage())
		return
	}
	response := response.SuccessResponse{
		Code: http.StatusText(http.StatusOK),
		Message: constant.OTP_VERIFICATION_SUCCEED,
	}
	utils.SendSuccessResponse(w, http.StatusOK, response)
}

// ResetUserPassword reset user's password to our database.
// @Summary reset user password
// @Description ResetUserPassword API is made for reset user password.
// @Accept json
// @Produce json
// @Tags users
// @Param id formData int64 true "ID which you've received in response of SendOTPToUser API"
// @Param password formData string true "New password of the user"
// @Success 200 {object} response.SuccessResponse "Password reset done successfully."
// @Failure 400 {object} errorhandling.CustomError "Bad request."
// @Failure 401 {object} errorhandling.CustomError "OTP not verified with our system."
// @Failure 500 {object} errorhandling.CustomError "Internal server error."
// @Router /api/v1/users/reset-password [put]
func (u userController) ResetUserPassword(w http.ResponseWriter, r *http.Request) {
	var userPasswordWithOTPId request.UserPasswordWithOTPID

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, errorhandling.ReadBodyError, constant.EMPTY_STRING)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &userPasswordWithOTPId)
	if err != nil {
		errorhandling.HandleJSONUnmarshlError(r, w, err)
		return
	}

	r.Body = io.NopCloser(bytes.NewReader(body))

	err = utils.Validate.Struct(userPasswordWithOTPId)
	if err != nil {
		errorhandling.HandleInvalidRequestData(w, r, err, utils.Translator)
		return
	}

	err = u.userService.ResetUserPassword(userPasswordWithOTPId)
	if err != nil {
		errorhandling.SendErrorResponse(r, w, err, utils.CreateErrorMessage())
		return
	}
	response := response.SuccessResponse{
		Message: "Password Reset Done Successfully.",
	}
	utils.SendSuccessResponse(w, http.StatusOK, response)
}
