package controller

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/chirag1807/task-management-system/api/model/dto"
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

// GetAllPublicProfileUsers fetches all public profile users.
// @Summary Get all public profile users
// @Description Get all public profile users based on query parameters
// @Produce json
// @Tags users
// @Param Authorization header string true "Access Token" default(Bearer <access_token>)
// @Param Limit query int false "Number of users to return per page (default 10)"
// @Param Offset query int false "Offset for pagination (default 0)"
// @Param Search query string false "Search term to filter users"
// @Success 200 {object} response.Users "Public profile users fetched successfully"
// @Failure 400 {object} errorhandling.CustomError "Bad request"
// @Failure 401 {object} errorhandling.CustomError "Either refresh token not found or token is expired."
// @Failure 500 {object} errorhandling.CustomError "Internal server error."
// @Router /api/user/get-public-profile-users [get]
func (u userController) GetAllPublicProfileUsers(w http.ResponseWriter, r *http.Request) {
	var queryParams = map[string]string{
		constant.LimitKey:          "number|default:10",
		constant.OffsetKey:         "number|default:0",
		constant.SearchKey:         "string",
	}
	var queryParamFilters = map[string]string{
		constant.LimitKey:          "int",
		constant.OffsetKey:         "int",
	}

	var userQueryParams request.UserQueryParams

	err, invalidParamsMultiLineErrMsg := utils.ValidateParameters(r, &userQueryParams, nil, nil, &queryParams, &queryParamFilters, nil)
	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}
	if invalidParamsMultiLineErrMsg != nil {
		errorhandling.SendErrorResponse(w, invalidParamsMultiLineErrMsg)
		return
	}

	publicProfileUsers, err := u.userService.GetAllPublicProfileUsers(userQueryParams)
	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}
	response := response.Users{
		Users: publicProfileUsers,
	}
	utils.SendSuccessResponse(w, http.StatusOK, response)
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
// @Router /api/user/get-my-details [get]
// /get-my-details
func (u userController) GetMyDetails(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(constant.UserIdKey).(int64)
	userDetails, err := u.userService.GetMyDetails(userId)
	if err != nil {
		errorhandling.SendErrorResponse(w, err)
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
// @Param profile formData string false "Profile of the user (Public, Private)"
// @Success 200 {object} response.SuccessResponse "User Updated successfully."
// @Failure 400 {object} errorhandling.CustomError "Bad request."
// @Failure 401 {object} errorhandling.CustomError "Either password not matched or need to left from all teams or token expired."
// @Failure 404 {object} errorhandling.CustomError "No user found."
// @Failure 409 {object} errorhandling.CustomError "Duplicate email found."
// @Failure 500 {object} errorhandling.CustomError "Internal server error."
// @Router /api/user/update-user-profile [put]
func (u userController) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	var requestParams = map[string]string{
		constant.FirstNameKey: "string|minLen:2",
		constant.LastNameKey:  "string|minLen:2",
		constant.BioKey:       "string|minLen:6",
		constant.EmailKey:     `string|regex:^[\w.%+-]+@[\w.-]+\.[a-zA-Z]{2,}$`,
		constant.PasswordKey:  "string|minLen:8",
		constant.ProfileKey:   "string|in:Public,Private",
	}
	var userToUpdate request.User

	err, invalidParamsMultiLineErrMsg := utils.ValidateParameters(r, &userToUpdate, &requestParams, nil, nil, nil, nil)
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

	err = json.Unmarshal(body, &userToUpdate)
	if err != nil {
		errorhandling.SendErrorResponse(w, errorhandling.ReadDataError)
		return
	}

	userId := r.Context().Value(constant.UserIdKey).(int64)
	if userToUpdate.Password != "" {
		err := u.userService.VerifyUserPassword(userToUpdate.Password, userId)
		if err != nil {
			errorhandling.SendErrorResponse(w, err)
			return
		}
		hashedPassword, err := utils.HashPassword(userToUpdate.NewPassword)
		if err != nil {
			errorhandling.SendErrorResponse(w, err)
			return
		}
		userToUpdate.Password = hashedPassword
	}

	err = u.userService.UpdateUserProfile(userId, userToUpdate)
	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}

	response := response.SuccessResponse{
		Message: constant.USER_PROFILE_UPDATED,
	}
	log.Println(constant.USER_PROFILE_UPDATED)
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
// @Router /api/user/send-otp-to-user [post]
func (u userController) SendOTPToUser(w http.ResponseWriter, r *http.Request) {
	var requestParams = map[string]string{
		constant.EmailKey: `string|regex:^[\w.%+-]+@[\w.-]+\.[a-zA-Z]{2,}$|required`,
	}
	var user request.User

	err, invalidParamsMultiLineErrMsg := utils.ValidateParameters(r, &user, &requestParams, nil, nil, nil, nil)
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

	err = json.Unmarshal(body, &user)
	if err != nil {
		errorhandling.SendErrorResponse(w, errorhandling.ReadDataError)
		return
	}

	min, max := 1001, 9999
	OTP := rand.Intn(max-min) + min

	emailBody := utils.PrepareEmailBody(OTP)
	email := dto.Email{
		To:      user.Email,
		Subject: "OTP Verification",
		Body:    emailBody,
	}

	otpExpireTime := time.Now().UTC().Add(time.Minute * 5)
	id, err := u.userService.SendOTPToUser(email, OTP, otpExpireTime)
	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}

	log.Println(OTP)
	response := response.SuccessResponse{
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
// @Failure 403 {object} errorhandling.CustomError "OTP verification time expired."
// @Failure 500 {object} errorhandling.CustomError "Internal server error."
// @Router /api/user/verify-otp [post]
func (u userController) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var requestParams = map[string]string{
		constant.OTPIdKey:   "number|required",
		constant.OTPCodeKey: "int|min:1000|max:9999|required",
	}
	var otp request.OTP

	err, invalidParamsMultiLineErrMsg := utils.ValidateParameters(r, &otp, &requestParams, nil, nil, nil, nil)
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

	err = json.Unmarshal(body, &otp)
	if err != nil {
		errorhandling.SendErrorResponse(w, errorhandling.ReadDataError)
		return
	}

	err = u.userService.VerifyOTP(otp)
	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}
	response := response.SuccessResponse{
		Message: "OTP Verification Done Successfully, You can proceed Further.",
	}
	utils.SendSuccessResponse(w, http.StatusOK, response)
}

// ResetUserPassword reset user's password to our database.
// @Summary reset user password
// @Description ResetUserPassword API is made for reset user password.
// @Accept json
// @Produce json
// @Tags users
// @Param email formData string true "Email of the user"
// @Param password formData string true "New password of the user"
// @Success 200 {object} response.SuccessResponse "Password reset done successfully."
// @Failure 400 {object} errorhandling.CustomError "Bad request."
// @Failure 500 {object} errorhandling.CustomError "Internal server error."
// @Router /api/user/reset-user-password [put]
func (u userController) ResetUserPassword(w http.ResponseWriter, r *http.Request) {
	var requestParams = map[string]string{
		constant.EmailKey:     `string|regex:^[\w.%+-]+@[\w.-]+\.[a-zA-Z]{2,}$|required`,
		constant.PasswordKey:  "string|minLen:8|required",
	}
	var userEmailPassword request.User

	err, invalidParamsMultiLineErrMsg := utils.ValidateParameters(r, &userEmailPassword, &requestParams, nil, nil, nil, nil)
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

	err = json.Unmarshal(body, &userEmailPassword)
	if err != nil {
		errorhandling.SendErrorResponse(w, errorhandling.ReadDataError)
		return
	}

	err = u.userService.ResetUserPassword(userEmailPassword)
	if err != nil {
		errorhandling.SendErrorResponse(w, err)
		return
	}
	response := response.SuccessResponse{
		Message: "Password Reset Done Successfully.",
	}
	utils.SendSuccessResponse(w, http.StatusOK, response)
}
