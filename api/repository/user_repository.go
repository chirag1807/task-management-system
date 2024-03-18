package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/chirag1807/task-management-system/api/model/dto"
	"github.com/chirag1807/task-management-system/api/model/request"
	"github.com/chirag1807/task-management-system/api/model/response"
	errorhandling "github.com/chirag1807/task-management-system/error"
	"github.com/chirag1807/task-management-system/utils"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	GetAllPublicProfileUsers() ([]response.User, error)
	GetMyDetails(userId int64) (response.User, error)
	UpdateUserProfile(userId int64, userToUpdate request.User) error
	SendOTPToUser(userEmail dto.Email, OTP int, OTPExpireTime time.Time) (int64, error)
	VerifyOTP(otpFromUser request.OTP) error
	ResetUserPassword(userEmailPassword request.User) error
	VerifyUserPassword(userPassword string, userId int64) error
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

func (u userRepository) GetAllPublicProfileUsers() ([]response.User, error) {
	publicUsers, err := u.dbConn.Query(context.Background(), `SELECT * FROM users WHERE profile = $1`, "Public")
	publicUsersSlice := make([]response.User, 0)
	if err != nil {
		return publicUsersSlice, err
	}
	defer publicUsers.Close()

	var publicUser response.User
	for publicUsers.Next() {
		if err := publicUsers.Scan(&publicUser.ID, &publicUser.FirstName, &publicUser.LastName, &publicUser.Bio, &publicUser.Email, &publicUser.Password, &publicUser.Profile); err != nil {
			return publicUsersSlice, err
		}
		publicUsersSlice = append(publicUsersSlice, publicUser)
	}
	return publicUsersSlice, nil
}

func (u userRepository) GetMyDetails(userId int64) (response.User, error) {
	var userDetails response.User
	user := u.dbConn.QueryRow(context.Background(), `SELECT * FROM users WHERE id = $1`, userId)
	err := user.Scan(&userDetails.ID, &userDetails.FirstName, &userDetails.LastName, &userDetails.Bio, &userDetails.Email, &userDetails.Password, &userDetails.Profile)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return userDetails, errorhandling.NoUserFound
		}
		return userDetails, err
	}
	return userDetails, nil
}

func (u userRepository) UpdateUserProfile(userId int64, userToUpdate request.User) error {
	query, args, err := UpdateQuery("users", userToUpdate, userId)
	if err != nil {
		return err
	}
	_, err = u.dbConn.Exec(context.Background(), query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (u userRepository) SendOTPToUser(userEmail dto.Email, OTP int, OTPExpireTime time.Time) (int64, error) {
	ctx := context.Background()
	var databaseOTPId int64
	var userCount int
	u.dbConn.QueryRow(ctx, `SELECT COUNT(*) FROM users where email = $1`, userEmail.To).Scan(&userCount)
	if userCount == 0 {
		return 0, errorhandling.NoEmailFound
	}
	tx, err := u.dbConn.Begin(ctx)
	if err != nil {
		return 0, err
	}
	err = tx.QueryRow(ctx, `INSERT INTO otps (otp, otp_expire_time) VALUES ($1, $2) RETURNING id`, OTP, OTPExpireTime).Scan(&databaseOTPId)
	if err != nil {
		tx.Rollback(ctx)
		return 0, err
	}
	err = utils.SendEmail(userEmail)
	if err != nil {
		tx.Rollback(ctx)
		return 0, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return 0, err
	}
	return databaseOTPId, nil
}

func (u userRepository) VerifyOTP(otpFromUser request.OTP) error {
	var dbOTP response.OTP
	rows, err := u.dbConn.Query(context.Background(), `SELECT * FROM otps where id = $1`, otpFromUser.ID)
	if err != nil {
		return err
	}
	if rows.Next() {
		if err := rows.Scan(&dbOTP.ID, &dbOTP.OTP, &dbOTP.OTPExpiryTime); err != nil {
			return err
		}
		fmt.Println(dbOTP.OTPExpiryTime.Before(time.Now()))
		fmt.Println(dbOTP.OTPExpiryTime.After(time.Now()))
		fmt.Println(dbOTP.OTPExpiryTime)
		fmt.Println(time.Now())
		if dbOTP.OTPExpiryTime.After(time.Now()) {
			return errorhandling.OTPVerificationTimeExpired
		} else if dbOTP.OTP != otpFromUser.OTP {
			return errorhandling.OTPNotMatched
		} else {
			return nil
		}
	} else {
		return errorhandling.NoOTPIDFound
	}
}

func (u userRepository) ResetUserPassword(userEmailPassword request.User) error {
	// var userCount int
	// u.dbConn.QueryRow(context.Background(), `SELECT COUNT(*) FROM users where email = $1`, userEmailPassword.Email).Scan(&userCount)
	// if userCount == 0 {
	// 	return errorhandling.NoEmailFound
	// }
	_, err := u.dbConn.Exec(context.Background(), "UPDATE users SET password = $1 WHERE email = $2", userEmailPassword.NewPassword, userEmailPassword.Email)
	if err != nil {
		return err
	}
	return nil
}

func (u userRepository) VerifyUserPassword(userPassword string, userId int64) error {
	var dbUser response.User
	row := u.dbConn.QueryRow(context.Background(), `SELECT password FROM users WHERE id = $1`, userId)
	err := row.Scan(&dbUser.Password)

	if err != nil && err.Error() == "no rows in result set" {
		return errorhandling.NoUserFound
	}

	passwordMatched := utils.VerifyPassword(userPassword, dbUser.Password)
	if !passwordMatched {
		return errorhandling.PasswordNotMatched
	}

	return nil
}
