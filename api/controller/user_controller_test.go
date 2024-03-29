package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/chirag1807/task-management-system/api/model/request"
	"github.com/chirag1807/task-management-system/constant"
	"github.com/stretchr/testify/assert"
)

func TestGetAllPublicProfileUsers(t *testing.T) {
	testCases := []struct {
		TestCaseName string
		QueryParams  request.UserQueryParams
		Expected     interface{}
		StatusCode   int
	}{
		{
			TestCaseName: "All Public Profile Users Fetched Successfully",
			QueryParams: request.UserQueryParams{
				Limit:  1,
				Offset: 0,
				Search: "",
			},
			StatusCode: 200,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {
			r.Get("/api/user/get-public-profile-users", NewUserController(userService).GetAllPublicProfileUsers)

			req, err := http.NewRequest("GET", "/api/user/get-public-profile-users", http.NoBody)
			if err != nil {
				log.Println(err)
			}

			q := req.URL.Query()
			q.Add("limit", strconv.Itoa(v.QueryParams.Limit))
			q.Add("offset", strconv.Itoa(v.QueryParams.Offset))
			q.Add("search", v.QueryParams.Search)
			req.URL.RawQuery = q.Encode()

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, v.StatusCode, w.Code)
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
			TestCaseName: "All Public Profile Users Fetched Successfully",
			UserID:       954488202459119617,
			StatusCode:   200,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {
			r.Get("/api/user/get-my-details", NewUserController(userService).GetMyDetails)

			req, err := http.NewRequest("GET", "/api/user/get-my-details", http.NoBody)
			if err != nil {
				log.Println(err)
			}

			ctx := context.WithValue(req.Context(), constant.UserIdKey, v.UserID)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, v.StatusCode, w.Code)
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
		Profile      string
		UserID       int64
		Expected     interface{}
		StatusCode   int
	}{
		{
			TestCaseName: "Profile Updated Successfully.",
			FirstName:    "Dhyey D",
			LastName:     "Panchal",
			UserID:       954488202459119617,
			Expected:     "Profile Updated Successfully.",
			StatusCode:   200,
		},
		{
			TestCaseName: "Field Must be of Minimum Length",
			FirstName:    "Dhyey",
			LastName:     "P",
			UserID:       954488202459119617,
			Expected:     "lastName violates minimum length constraint.",
			StatusCode:   400,
		},
		{
			TestCaseName: "Invalid Email",
			Email:        "dhyeypanchal2204",
			UserID:       954488202459119617,
			Expected:     "please provide email in valid format.",
			StatusCode:   400,
		},
		{
			TestCaseName: "Duplicate Email",
			Email:        "ridham@gmail.com",
			UserID:       954488202459119617,
			Expected:     "Duplicate Email Found.",
			StatusCode:   409,
		},
		{
			TestCaseName: "Value Must be in Enum Values.",
			Profile:      "public",
			UserID:       954488202459119617,
			Expected:     "profile value must be in the enum [Public Private]",
			StatusCode:   400,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {
			r.Post("/api/user/update-user-profile", NewUserController(userService).UpdateUserProfile)

			user := request.User{
				FirstName: v.FirstName,
				LastName:  v.LastName,
				Bio:       v.Bio,
				Email:     v.Email,
				Password:  v.Password,
				Profile:   v.Profile,
			}
			jsonValue, err := json.Marshal(user)
			if err != nil {
				log.Println(err)
			}
			req, err := http.NewRequest("POST", "/api/user/update-user-profile", bytes.NewBuffer(jsonValue))
			if err != nil {
				log.Println(err)
			}
			req.Header.Set("Content-Type", "application/json")

			req.Header.Set("Content-Type", "application/json")
			ctx := context.WithValue(req.Context(), constant.UserIdKey, v.UserID)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, v.StatusCode, w.Code)
		})
	}
}

func TestSendOTPToUser(t *testing.T) {
	testCases := []struct {
		TestCaseName string
		Email        string
		Expected     interface{}
		StatusCode   int
	}{
		{
			TestCaseName: "OTP Sent Successfully.",
			Email:        "dhyey@gmail.com",
			Expected:     "User Registration Done Successfully.",
			StatusCode:   200,
		},
		{
			TestCaseName: "Field Must be Required",
			Expected:     "Email is required field.",
			StatusCode:   400,
		},
		{
			TestCaseName: "Invalid Email",
			Email:        "dhyeypanchal2204",
			Expected:     "please provide email in valid format.",
			StatusCode:   400,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {
			r.Post("/api/user/send-otp-to-user", NewUserController(userService).SendOTPToUser)

			user := request.User{
				Email: v.Email,
			}
			jsonValue, err := json.Marshal(user)
			if err != nil {
				log.Println(err)
			}
			req, err := http.NewRequest("POST", "/api/user/send-otp-to-user", bytes.NewBuffer(jsonValue))
			if err != nil {
				log.Println(err)
			}
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, v.StatusCode, w.Code)
		})
	}
}

func TestVerifyOTP(t *testing.T) {
	testCases := []struct {
		TestCaseName string
		OTPID        int64
		OTP          int
		Expected     interface{}
		StatusCode   int
	}{
		{
			TestCaseName: "OTP Verified Successfully.",
			OTPID:        954537852771565569,
			OTP:          1099,
			StatusCode:   200,
		},
		{
			TestCaseName: "Field Must be Required",
			OTP:          1,
			StatusCode:   400,
		},
		{
			TestCaseName: "Field Violates Minimum Value Constraint",
			OTPID:        1,
			OTP:          1,
			StatusCode:   400,
		},
		{
			TestCaseName: "Field Violates Maximum Value Constraint",
			OTPID:        1,
			OTP:          10000,
			StatusCode:   400,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {
			r.Post("/api/user/verify-otp", NewUserController(userService).VerifyOTP)

			user := request.OTP{
				ID:  v.OTPID,
				OTP: v.OTP,
			}
			jsonValue, err := json.Marshal(user)
			if err != nil {
				log.Println(err)
			}
			req, err := http.NewRequest("POST", "/api/user/verify-otp", bytes.NewBuffer(jsonValue))
			if err != nil {
				log.Println(err)
			}
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, v.StatusCode, w.Code)
		})
	}
}

func TestResetUserPassword(t *testing.T) {
	testCases := []struct {
		TestCaseName string
		Email        string
		Password     string
		Expected     interface{}
		StatusCode   int
	}{
		{
			TestCaseName: "Password Reset Done Successfully.",
			Email:        "dhyey@gmail.com",
			Password:     "Dhyey123$",
			Expected:     "User Registration Done Successfully.",
			StatusCode:   200,
		},
		{
			TestCaseName: "Field Must be Required",
			Email:        "dhyey@gmail.com",
			Expected:     "password is required field.",
			StatusCode:   400,
		},
		{
			TestCaseName: "Field Must be of Minimum Length",
			Email:        "dhyey@gmail.com",
			Password:     "Dhyey",
			Expected:     "lastName violates minimum length constraint.",
			StatusCode:   400,
		},
		{
			TestCaseName: "Invalid Email",
			Email:        "dhyeypanchal2204",
			Password:     "Dhyey123$",
			Expected:     "please provide email in valid format.",
			StatusCode:   400,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {
			r.Put("/api/user/reset-user-password", NewUserController(userService).ResetUserPassword)

			user := request.User{
				Email:    v.Email,
				Password: v.Password,
			}
			jsonValue, err := json.Marshal(user)
			if err != nil {
				log.Println(err)
			}
			req, err := http.NewRequest("PUT", "/api/user/reset-user-password", bytes.NewBuffer(jsonValue))
			if err != nil {
				log.Println(err)
			}
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, v.StatusCode, w.Code)
		})
	}
}
