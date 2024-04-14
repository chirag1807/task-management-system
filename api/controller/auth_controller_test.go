package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chirag1807/task-management-system/api/model/request"
	"github.com/chirag1807/task-management-system/constant"
	"github.com/stretchr/testify/assert"
)

func TestUserRegistration(t *testing.T) {
	testCases := []struct {
		TestCaseName string
		FirstName    string
		LastName     string
		Bio          string
		Email        string
		Password     string
		Privacy      string
		Expected     interface{}
		StatusCode   int
	}{
		{
			TestCaseName: "Registration Done Successfully.",
			FirstName:    "Chirag",
			LastName:     "Makwana",
			Bio:          "Junior Software Engineer",
			Email:        "chiragmakwana1807@gmail.com",
			Password:     "Chirag123$",
			Privacy:      "PUBLIC",
			Expected:     "User Registration Done Successfully.",
			StatusCode:   200,
		},
		{
			TestCaseName: "Field Must be Not Empty.",
			FirstName:    "Chirag",
			LastName:     "",
			Bio:          "Junior Software Engineer",
			Email:        "chiragmakwana1807@gmail.com",
			Password:     "Chirag123$",
			Privacy:      "PUBLIC",
			Expected:     "lastName is required to not be empty.",
			StatusCode:   400,
		},
		{
			TestCaseName: "Field Must be Required",
			FirstName:    "Chirag",
			Bio:          "Junior Software Engineer",
			Email:        "chiragmakwana1807@gmail.com",
			Password:     "Chirag123$",
			Privacy:      "PUBLIC",
			Expected:     "lastName is required to not be empty.",
			StatusCode:   400,
		},
		{
			TestCaseName: "Field Must be of Minimum Length",
			FirstName:    "Chirag",
			LastName:     "M",
			Bio:          "Junior Software Engineer",
			Email:        "chiragmakwana1807@gmail.com",
			Password:     "Chirag123$",
			Privacy:      "PUBLIC",
			Expected:     "lastName violates minimum length constraint.",
			StatusCode:   400,
		},
		{
			TestCaseName: "Invalid Email",
			FirstName:    "Chirag",
			LastName:     "Makwana",
			Bio:          "Junior Software Engineer",
			Email:        "chiragmakwana1807",
			Password:     "Chirag123$",
			Privacy:      "PUBLIC",
			Expected:     "please provide email in valid format.",
			StatusCode:   400,
		},
		{
			TestCaseName: "Duplicate Email",
			FirstName:    "Chirag",
			LastName:     "Makwana",
			Bio:          "Junior Software Engineer",
			Email:        "chiragmakwana1807@gmail.com",
			Password:     "Chirag123$",
			Privacy:      "PUBLIC",
			Expected:     "Duplicate Email Found.",
			StatusCode:   409,
		},
		{
			TestCaseName: "Value Must be in Enum Values.",
			FirstName:    "Chirag",
			LastName:     "Makwana",
			Bio:          "Junior Software Engineer",
			Email:        "chiragmakwana18@gmail.com",
			Password:     "Chirag123$",
			Privacy:      "public",
			Expected:     "privacy value must be in the enum [PUBLIC PRIVATE]",
			StatusCode:   400,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {
			r.Post("/api/v1/auth/registration", NewAuthController(authService).UserRegistration)

			user := request.User{
				FirstName: v.FirstName,
				LastName:  v.LastName,
				Bio:       v.Bio,
				Email:     v.Email,
				Password:  v.Password,
				Privacy:   v.Privacy,
			}
			jsonValue, err := json.Marshal(user)
			if err != nil {
				
				log.Println(err)
			}
			req, err := http.NewRequest("POST", "/api/v1/auth/registration", bytes.NewBuffer(jsonValue))
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

func TestUserLogin(t *testing.T) {
	testCases := []struct {
		TestCaseName string
		Email        string
		Password     string
		StatusCode   int
	}{
		{
			TestCaseName: "Login Done Successfully.",
			Email:        "guptaaahutosh354@gmail.com",
			Password:     "Aashutosh1234$",
			StatusCode:   200,
		},
		{
			TestCaseName: "Field Must be Not Empty.",
			Email:        "chiragmakwana1807@gmail.com",
			Password:     "",
			StatusCode:   400,
		},
		{
			TestCaseName: "Field Must be Required",
			Email:        "guptaaahutosh354@gmail.com",
			StatusCode:   400,
		},
		{
			TestCaseName: "Invalid Email",
			Email:        "guptaaahutosh354@gmail",
			Password:     "Aashutosh1234$",
			StatusCode:   400,
		},
		{
			TestCaseName: "Email Not Found",
			Email:        "nirajdarji@gmail.com",
			Password:     "Aashutosh1234$",
			StatusCode:   404,
		},
		{
			TestCaseName: "Wrong Password",
			Email:        "guptaaahutosh354@gmail.com",
			Password:     "Chirag123$",
			StatusCode:   401,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {
			r.Post("/api/v1/auth/login", NewAuthController(authService).UserLogin)

			user := request.User{
				Email:    v.Email,
				Password: v.Password,
			}
			jsonValue, _ := json.Marshal(user)
			req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, v.StatusCode, w.Code)
		})
	}
}

func TestResetToken(t *testing.T) {
	testCases := []struct {
		TestCaseName string
		Token        string
		StatusCode   int
	}{
		{
			TestCaseName: "Token Reset Done Successfully.",
			Token:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTE3OTMxMjUsInVzZXJJZCI6Ijk1MzkzNDU1MzI1NDEwMDk5MyJ9.vFcrOMncN7y8nBkWV6iULeafZLp73z7kNZDzb2e0-PM",
			StatusCode:   200,
		},
		{
			TestCaseName: "Refresh Token Not Valid or Expired.",
			Token:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			StatusCode:   401,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {
			r.Post("/api/v1/auth/refresh-token", NewAuthController(authService).RefreshToken)

			req, _ := http.NewRequest("POST", "/api/v1/auth/refresh-token", http.NoBody)
			ctx := context.WithValue(req.Context(), constant.TokenKey, v.Token)
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, v.StatusCode, w.Code)
		})
	}
}
