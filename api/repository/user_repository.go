package repository

import (
	"github.com/chirag1807/task-management-system/api/model/request"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	UpdateUserProfile(userToUpdate request.User) error
	SendOTPToUser(userEmail string) (int8, error)
	VerifyOTP(otpFromUser request.OTP) error
	ResetUserPassword(userEmailPassword request.User) error
}

type userRepository struct {
	dbConn      *pgx.Conn
	redisClient *redis.Client
}

func NewUserRepo(dbConn *pgx.Conn, redisClient *redis.Client) UserRepository {
	return userRepository{
		dbConn:      dbConn,
		redisClient: redisClient,
	}
}

func (u userRepository) UpdateUserProfile(userToUpdate request.User) error {
	return nil
}

func (u userRepository) SendOTPToUser(userEmail string) (int8, error) {
	return 1, nil
}

func (u userRepository) VerifyOTP(otpFromUser request.OTP) error {
	return nil
}

func (u userRepository) ResetUserPassword(userEmailPassword request.User) error {
	return nil
}
