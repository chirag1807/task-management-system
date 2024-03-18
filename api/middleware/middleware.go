package middleware

import (
	"context"
	"net/http"

	"github.com/chirag1807/task-management-system/constant"
	errorhandling "github.com/chirag1807/task-management-system/error"
	"github.com/chirag1807/task-management-system/utils"
)

// VerifyToken retrieves token from request header and send it to VerifyJWTToken function of utils package.
// after that it will check that err is nil or not and if it is nil then send token expired error response from here.
// otherwise it will set token and userId to request's context and command will go to controller section.
func VerifyToken(flag int) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				errorhandling.SendErrorResponse(w, errorhandling.TokenNotFound)
				return
			}
			token = token[7:]
			userId, err := utils.VerifyJWTToken(token)
			if err != nil {
				if flag == 0 {
					errorhandling.SendErrorResponse(w, errorhandling.AccessTokenExpired)
				} else {
					errorhandling.SendErrorResponse(w, errorhandling.RefreshTokenExpired)
				}
				return
			}
			ctx := context.WithValue(r.Context(), constant.TokenKey, token)
			ctx = context.WithValue(ctx, constant.UserIdKey, userId)
			handler.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
