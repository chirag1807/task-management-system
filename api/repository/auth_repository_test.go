package repository

import (
	"testing"

	"github.com/chirag1807/task-management-system/api/model/request"
	errorhandling "github.com/chirag1807/task-management-system/error"
	"github.com/stretchr/testify/assert"
)

func TestRefreshToken(t *testing.T) {
	testCases := []struct {
		TestCaseName string
		Token        string
		Expected     interface{}
	}{
		{
			TestCaseName: "Token Rest Done Successfully.",
			Token:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ2OTE0OTMsInVzZXJJZCI6Ijk1NDQ4ODIwMjQ1OTExOTYxNyJ9.qi3BFn6UhmodlODzSNfGVxzLxjsCncM7GPvVZya5aLc",
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
			_, _, err := NewAuthRepo(dbConn).RefreshToken(v.Token)

			assert.Equal(t, v.Expected, err)
		})
	}
}

func TestUserRegistration(t *testing.T) {
	testCases := []struct {
		TestCaseName string
		FirstName    string
		LastName     string
		Bio          string
		Email        string
		Password     string
		Privacy      string
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
			Privacy:      "PUBLIC",
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
			Privacy:      "PUBLIC",
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
				Privacy:   v.Privacy,
			}

			_, err := NewAuthRepo(dbConn).UserRegistration(user)
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
			user := request.UserCredentials{
				Email:    v.Email,
				Password: v.Password,
			}
			_, _, err := NewAuthRepo(dbConn).UserLogin(user)

			assert.Equal(t, v.Expected, err)
		})
	}
}
