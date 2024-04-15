package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/chirag1807/task-management-system/api/model/request"
	"github.com/chirag1807/task-management-system/api/model/response"
	"github.com/chirag1807/task-management-system/constant"
	errorhandling "github.com/chirag1807/task-management-system/error"
	"github.com/chirag1807/task-management-system/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type AuthRepository interface {
	UserRegistration(user request.User) (int64, error)
	UserLogin(user request.UserCredentials) (response.User, string, error)
	RefreshToken(refreshToken string) (int64, string, error)
}

type authRepository struct {
	dbConn *pgx.Conn
}

func NewAuthRepo(dbConn *pgx.Conn) AuthRepository {
	return authRepository{
		dbConn: dbConn,
	}
}

func (a authRepository) UserRegistration(user request.User) (int64, error) {
	var userID int64
	fmt.Println(user.Email)
	rows := a.dbConn.QueryRow(context.Background(), `INSERT INTO users (first_name, last_name, bio, email, password, privacy) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, user.FirstName, user.LastName, user.Bio, user.Email, user.Password, user.Privacy)
	err := rows.Scan(&userID)
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if ok && pgErr.Code == constant.PG_Duplicate_Error_Code {
			return 0, errorhandling.DuplicateEmailFound
		}
		return 0, err
	}
	return userID, nil
}

func (a authRepository) UserLogin(user request.UserCredentials) (response.User, string, error) {
	var dbUser response.User
	rows := a.dbConn.QueryRow(context.Background(), `SELECT id, first_name, last_name, bio, email, password, privacy FROM users WHERE email = $1`, user.Email)
	err := rows.Scan(&dbUser.ID, &dbUser.FirstName, &dbUser.LastName, &dbUser.Bio, &dbUser.Email, &dbUser.Password, &dbUser.Privacy)

	if err != nil && err.Error() == constant.PG_NO_ROWS {
		return response.User{}, constant.EMPTY_STRING, errorhandling.NoUserFound
	}

	passwordMatched := utils.VerifyPassword(user.Password, dbUser.Password)
	if !passwordMatched {
		return response.User{}, constant.EMPTY_STRING, errorhandling.PasswordNotMatched
	}

	refreshToken, err := utils.CreateJWTToken(time.Now().Add(time.Hour*24*7), dbUser.ID)
	if err != nil {
		return response.User{}, constant.EMPTY_STRING, err
	}

	_, err = a.dbConn.Exec(context.Background(), `INSERT INTO refresh_tokens (user_id, refresh_token) VALUES ($1, $2)`, dbUser.ID, refreshToken)
	if err != nil {
		return response.User{}, constant.EMPTY_STRING, err
	}

	return dbUser, refreshToken, nil
}

func (a authRepository) RefreshToken(refreshToken string) (int64, string, error) {
	var userID int64
	rows := a.dbConn.QueryRow(context.Background(), `SELECT id FROM users as u LEFT JOIN refresh_tokens as r on u.id = r.user_id WHERE r.refresh_token = $1`, refreshToken)
	err := rows.Scan(&userID)
	if err != nil {
		if err.Error() == constant.PG_NO_ROWS {
			return 0, constant.EMPTY_STRING, errorhandling.RefreshTokenNotFound
		} else {
			return 0, constant.EMPTY_STRING, err
		}
	}

	refreshToken, err = utils.CreateJWTToken(time.Now().Add(time.Hour*24*7), userID)
	if err != nil {
		return 0, constant.EMPTY_STRING, err
	}

	_, err = a.dbConn.Exec(context.Background(), `INSERT INTO refresh_tokens (user_id, refresh_token) VALUES ($1, $2)`, userID, refreshToken)
	if err != nil {
		return 0, constant.EMPTY_STRING, err
	}

	return userID, refreshToken, nil
}
