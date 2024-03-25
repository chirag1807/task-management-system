package repository

import (
	"testing"

	"github.com/chirag1807/task-management-system/api/model/request"
	errorhandling "github.com/chirag1807/task-management-system/error"
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
		Profile      string
		Expected     error
		StatusCode   int
	}{
		{
			TestCaseName: "Registration Done Successfully.",
			FirstName:    "Chirag",
			LastName:     "Makwana",
			Bio:          "Junior Software Engineer",
			Email:        "chiragmakwana1807@gmail.com",
			Password:     "Chirag123$",
			Profile:      "Public",
			Expected:     nil,
			StatusCode:   200,
		},
		{
			TestCaseName: "Duplicate Email",
			FirstName:    "Chirag",
			LastName:     "Makwana",
			Bio:          "Junior Software Engineer",
			Email:        "chiragmakwana1807@gmail.com",
			Password:     "Chirag123$",
			Profile:      "Public",
			Expected:     errorhandling.DuplicateEmailFound,
			StatusCode:   409,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {

			user := request.User{
				FirstName: v.FirstName,
				LastName:  v.LastName,
				Bio:       v.Bio,
				Email:     v.Email,
				Password:  v.Password,
				Profile:   v.Profile,
			}

			_, err := NewAuthRepo(dbConn, redisClient).UserRegistration(user)
			assert.Equal(t, v.Expected, err)
		})
	}
}

func TestUserLogin(t *testing.T) {
	testCases := []struct {
		TestCaseName string
		Email        string
		Password     string
		Expected     interface{}
	}{
		{
			TestCaseName: "Login Done Successfully",
			Email:        "guptaaahutosh354@gmail.com",
			Password:     "Aashutosh1234$",
			Expected:     nil,
		},
		{
			TestCaseName: "Password Not Matched",
			Email:        "guptaaahutosh354@gmail.com",
			Password:     "Aashutosh",
			Expected:     errorhandling.PasswordNotMatched,
		},
		{
			TestCaseName: "No User",
			Email:        "nirajdarji",
			Password:     "Niraj123$",
			Expected:     errorhandling.NoUserFound,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {
			user := request.User{
				Email:    v.Email,
				Password: v.Password,
			}
			_, _, err := NewAuthRepo(dbConn, redisClient).UserLogin(user)

			assert.Equal(t, v.Expected, err)
		})
	}
}

func TestResetToken(t *testing.T) {
	testCases := []struct {
		TestCaseName string
		Token        string
		Expected     interface{}
	}{
		{
			TestCaseName: "Token Rest Done Successfully.",
			Token:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTE3OTMxMjUsInVzZXJJZCI6Ijk1MzkzNDU1MzI1NDEwMDk5MyJ9.vFcrOMncN7y8nBkWV6iULeafZLp73z7kNZDzb2e0-PM",
			Expected:     nil,
		},
		{
			TestCaseName: "Refresh Token Error",
			Token:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			Expected:     errorhandling.RefreshTokenNotFound,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {
			_, err := NewAuthRepo(dbConn, redisClient).ResetToken(v.Token)

			assert.Equal(t, v.Expected, err)
		})
	}
}