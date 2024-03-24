package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chirag1807/task-management-system/api/model/request"
	"github.com/chirag1807/task-management-system/api/model/dto"
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
		Expected     interface{}
		StatusCode   int
	}{
		{
			TestCaseName: "Success",
			FirstName:    "Chirag",
			LastName:     "Makwana",
			Bio:          "Junior Software Engineer",
			Email:        "chiragmakwana1807@gmail.com",
			Password:     "Chirag123$",
			Profile:      "Public",
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
			Profile:      "Public",
			Expected:     "lastName is required to not be empty.",
			StatusCode:   400,
		},
		{
			TestCaseName: "Field Must be Required",
			FirstName:    "Chirag",
			Bio:          "Junior Software Engineer",
			Email:        "chiragmakwana1807@gmail.com",
			Password:     "Chirag123$",
			Profile:      "Public",
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
			Profile:      "Public",
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
			Profile:      "Public",
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
			Profile:      "Public",
			Expected:     "Duplicate Email Found.",
			StatusCode:   400,
		},
		{
			TestCaseName: "Value Must be in Enum Values.",
			FirstName:    "Chirag",
			LastName:     "Makwana",
			Bio:          "Junior Software Engineer",
			Email:        "chiragmakwana18@gmail.com",
			Password:     "Chirag123$",
			Profile:      "Public",
			Expected:     "profile value must be in the enum [Public Private]",
			StatusCode:   400,
		},
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {
			r.Post("/api/auth/registration", NewAuthController(authService).UserRegistration)

			user := request.User{
				FirstName: v.FirstName,
				LastName:  v.LastName,
				Bio:       v.Bio,
				Email:     v.Email,
				Password:  v.Password,
				Profile:   v.Profile,
			}
			jsonValue, _ := json.Marshal(user)
			req, _ := http.NewRequest("POST", "/api/auth/registration", bytes.NewBuffer(jsonValue))

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			var response dto.ExpectedMessage
			json.Unmarshal(w.Body.Bytes(), &response)

			assert.Equal(t, v.Expected, response.Message)
			assert.Equal(t, v.StatusCode, w.Code)
		})
	}
}
