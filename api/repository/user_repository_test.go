package repository

import (
	"testing"
	"time"

	"github.com/chirag1807/task-management-system/api/model/dto"
	"github.com/chirag1807/task-management-system/api/model/request"
	errorhandling "github.com/chirag1807/task-management-system/error"
	"github.com/chirag1807/task-management-system/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetAllPublicPrivacyUsers(t *testing.T) {
	testCases := []struct {
		TestCaseName string
		QueryParams  request.UserQueryParams
		Expected     interface{}
		StatusCode   int
	}{
		{
			TestCaseName: "Public Privacy Users Fetched Successfully",
			QueryParams: request.UserQueryParams{
				Limit:  1,
				Offset: 0,
				Search: "",
			},
			Expected:   nil,
			StatusCode: 200,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {
			_, err := NewUserRepo(dbConn, rabbitmqConn).GetAllPublicPrivacyUsers(v.QueryParams)
			assert.Equal(t, v.Expected, err)
		})
	}
}

func TestGetMyDetails(t *testing.T) {
	testCases := []struct {
		TestCaseName string
		UserID       int64
		Expected     interface{}
		StatusCode   int
	}{
		{
			TestCaseName: "Details Fetched Successfully",
			UserID:       954488202459119617,
			Expected:     nil,
			StatusCode:   200,
		},
		{
			TestCaseName: "Details Not Found for Not Existing Member",
			UserID:       954488202459119618,
			Expected:     errorhandling.NoUserFound,
			StatusCode:   200,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {
			_, err := NewUserRepo(dbConn, rabbitmqConn).GetMyDetails(v.UserID)
			assert.Equal(t, v.Expected, err)
		})
	}
}

func TestUpdateUserProfile(t *testing.T) {
	testCases := []struct {
		TestCaseName string
		FirstName    string
		LastName     string
		Bio          string
		Email        string
		Password     string
		Privacy      string
		UserID       int64
		Expected     interface{}
		StatusCode   int
	}{
		{
			TestCaseName: "Profile Updated Successfully.",
			FirstName:    "Dhyey Devendra Kumar",
			LastName:     "Panchal",
			UserID:       954488202459119617,
			Expected:     nil,
			StatusCode:   200,
		},
		{
			TestCaseName: "Leave All Teams to Make Privacy Private",
			Privacy:      "PRIVATE",
			UserID:       954488202459119617,
			Expected:     errorhandling.LeftAllTeamsToMakePrivacyPrivate,
			StatusCode:   401,
		},
		{
			TestCaseName: "Duplicate Email",
			Email:        "ridham@gmail.com",
			UserID:       954488202459119617,
			Expected:     errorhandling.DuplicateEmailFound,
			StatusCode:   409,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {
			userToUpdate := request.UpdateUser{
				FirstName: v.FirstName,
				LastName:  v.LastName,
				Email:     v.Email,
				Privacy:   v.Privacy,
			}
			err := NewUserRepo(dbConn, rabbitmqConn).UpdateUserProfile(v.UserID, userToUpdate)
			assert.Equal(t, v.Expected, err)
		})
	}
}

func TestSendOTPToUser(t *testing.T) {
	testCases := []struct {
		TestCaseName string
		Email        string
		OTP          int
		Expected     interface{}
		StatusCode   int
	}{
		{
			TestCaseName: "OTP Sent Successfully.",
			Email:        "dhyey@gmail.com",
			OTP:          1099,
			Expected:     nil,
			StatusCode:   200,
		},
		{
			TestCaseName: "No Email Found",
			Email:        "dhyey123@gmail.com",
			OTP:          1099,
			Expected:     errorhandling.NoEmailFound,
			StatusCode:   404,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {
			emailBody := utils.PrepareEmailBody(v.OTP)
			email := dto.Email{
				To:      v.Email,
				Subject: "OTP Verification",
				Body:    emailBody,
			}
			_, err := NewUserRepo(dbConn, rabbitmqConn).SendOTPToUser(email, 4099, time.Now().Add(5*time.Minute))
			assert.Equal(t, v.Expected, err)
		})
	}
}

func TestVerifyOTP(t *testing.T) {
	testCases := []struct {
		TestCaseName string
		ID           int64
		OTP          int
		Expected     interface{}
		StatusCode   int
	}{
		{
			TestCaseName: "OTP Sent Successfully.",
			OTP:          1099,
			ID:           954537852771565569,
			Expected:     nil,
			StatusCode:   200,
		},
		{
			TestCaseName: "No Email Found",
			ID:           954537852771565570,
			OTP:          1099,
			Expected:     errorhandling.NoOTPIDFound,
			StatusCode:   404,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {
			otpFromUser := request.OTP{
				ID:  v.ID,
				OTP: v.OTP,
			}
			err := NewUserRepo(dbConn, rabbitmqConn).VerifyOTP(otpFromUser)
			assert.Equal(t, v.Expected, err)
		})
	}
}

func TestResetUserPassword(t *testing.T) {
	testCases := []struct {
		TestCaseName string
		ID           int64
		Password     string
		Expected     interface{}
		StatusCode   int
	}{
		{
			TestCaseName: "Password Reset Done Successfully.",
			ID:           954537852771565569,
			Password:     "Dhyey123$",
			Expected:     nil,
			StatusCode:   200,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {
			userPasswordWithOTPId := request.UserPasswordWithOTPID{
				ID:       v.ID,
				Password: v.Password,
			}
			err := NewUserRepo(dbConn, rabbitmqConn).ResetUserPassword(userPasswordWithOTPId)
			assert.Equal(t, v.Expected, err)
		})
	}
}
