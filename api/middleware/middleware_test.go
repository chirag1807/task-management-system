package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	errorhandling "github.com/chirag1807/task-management-system/error"
	"github.com/magiconair/properties/assert"
)

func TestVerifyToken(t *testing.T) {
	testCases := []struct {
		TestCaseName string
		Token        string
		Flag         int
		Expected     interface{}
	}{
		{
			TestCaseName: "Valid Access Token",
			Token:        "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ2OTE0OTMsInVzZXJJZCI6Ijk1NDQ4ODIwMjQ1OTExOTYxNyJ9.qi3BFn6UhmodlODzSNfGVxzLxjsCncM7GPvVZya5aLc",
			Flag:         0,
			Expected:     nil,
		},
		{
			TestCaseName: "Valid Refresh Token",
			Token:        "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ2OTE0OTMsInVzZXJJZCI6Ijk1NDQ4ODIwMjQ1OTExOTYxNyJ9.qi3BFn6UhmodlODzSNfGVxzLxjsCncM7GPvVZya5aLc",
			Flag:         1,
			Expected:     nil,
		},
		// {
		// 	TestCaseName: "Token Not Found",
		// 	Flag:         1,
		// 	Expected:     errorhandling.TokenNotFound,
		// },
	}

	for _, v := range testCases {
		t.Run(v.TestCaseName, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/", nil)
			req.Header.Set("Authorization", v.Token)

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Do nothing for this test
			})

			rr := httptest.NewRecorder()
			VerifyToken(v.Flag)(handler).ServeHTTP(rr, req)

			if v.Expected != nil {
				assert.Equal(t, v.Expected, errorhandling.TokenNotFound)
			}
		})
	}
}
