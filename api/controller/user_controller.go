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
	utils.SendSuccessResponse(w, http.StatusOK, response)
}

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

	// emailBody, err := utils.ParseTemplate("D:/Task Manager GOLang/utils/email_body.html", map[string]string{"OTP": string(rune(OTP))})
	// if err != nil {
	// 	errorhandling.SendErrorResponse(w, err)
	// 	return
	// }

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
