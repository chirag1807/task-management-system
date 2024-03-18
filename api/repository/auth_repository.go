package repository

import (
	"context"
	"time"

	"github.com/chirag1807/task-management-system/api/model/request"
	"github.com/chirag1807/task-management-system/api/model/response"
	errorhandling "github.com/chirag1807/task-management-system/error"
	"github.com/chirag1807/task-management-system/utils"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type AuthRepository interface {
	UserRegistration(user request.User) (int64, error)
	UserLogin(user request.User) (response.User, string, error)
	ResetToken(refreshToken string) (int64, error)
}

type authRepository struct {
	dbConn      *pgx.Conn
	redisClient *redis.Client
}

func NewAuthRepo(dbConn *pgx.Conn, redisClient *redis.Client) AuthRepository {
	return authRepository{
		dbConn:      dbConn,
		redisClient: redisClient,
	}
}

func (a authRepository) UserRegistration(user request.User) (int64, error) {
	var userID int64
	err := a.dbConn.QueryRow(context.Background(), `INSERT INTO users (first_name, last_name, bio, email, password, profile) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, user.FirstName, user.LastName, user.Bio, user.Email, user.Password, user.Profile).Scan(&userID)
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if ok && pgErr.Code == "23505" {
			return 0, errorhandling.DuplicateEmailFound
		}
		return 0, err
	}
	return userID, nil
}

func (a authRepository) UserLogin(user request.User) (response.User, string, error) {
	var dbUser response.User
	row := a.dbConn.QueryRow(context.Background(), `SELECT id, first_name, last_name, bio, email, password, profile FROM users WHERE email = $1`, user.Email)
	err := row.Scan(&dbUser.ID, &dbUser.FirstName, &dbUser.LastName, &dbUser.Bio, &dbUser.Email, &dbUser.Password, &dbUser.Profile)

	if err != nil && err.Error() == "no rows in result set" {
		return response.User{}, "", errorhandling.NoUserFound
	}

	passwordMatched := utils.VerifyPassword(user.Password, dbUser.Password)
	if !passwordMatched {
		return response.User{}, "", errorhandling.PasswordNotMatched
	}

	refreshToken, err := utils.CreateJWTToken(time.Now().Add(time.Hour*24*7), dbUser.ID)
	if err != nil {
		return response.User{}, "", err
	}

	_, err = a.dbConn.Exec(context.Background(), `INSERT INTO refresh_tokens (user_id, refresh_token) VALUES ($1, $2)`, dbUser.ID, refreshToken)
	if err != nil {
		return response.User{}, "", err
	}

	return dbUser, refreshToken, nil
}

func (a authRepository) ResetToken(refreshToken string) (int64, error) {
	row := a.dbConn.QueryRow(context.Background(), `SELECT id FROM users as u LEFT JOIN refresh_tokens as r on u.id = r.user_id WHERE r.refresh_token = $1`, refreshToken)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return 0, errorhandling.RefreshTokenNotFound
		} else {
			return 0, err
		}
	}

	return id, nil
}
