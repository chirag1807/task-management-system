package repository

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/chirag1807/task-management-system/api/model/request"
	"github.com/chirag1807/task-management-system/api/model/response"
	errorhandling "github.com/chirag1807/task-management-system/error"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	GetAllPublicProfileUsers() ([]response.User, error)
	GetMyDetails(userId int64) (response.User, error)
	UpdateUserProfile(userId int64, userToUpdate request.UpdateUser) error
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

func (u userRepository) UpdateUserProfile(userId int64, userToUpdate request.UpdateUser) error {

	var (
		columns []string
		values  []interface{}
	)

	map_2 := map[string]string{
		"FirstName": "first_name",
		"LastName":  "last_name",
		"Bio":       "bio",
		"Email":     "email",
		"Password":  "password",
		"Profile":   "profile",
	}
	counter := 1

	fmt.Println(*(userToUpdate.Profile))
	v := reflect.ValueOf(userToUpdate)
	typeOf := v.Type()
	for i := 0; i < v.NumField(); i++ {
		if !(v.Field(i).IsNil()) {
			columns = append(columns, fmt.Sprintf("%s = $%v,", map_2[typeOf.Field(i).Name], counter))
			counter++
			values = append(values, v.Field(i).Interface())
		}
	}

	queryStr := "UPDATE users SET " + strings.Join(columns, " ") + "WHERE id = $" + fmt.Sprint(counter)
	fmt.Println(queryStr)
	log.Println(values...)

	// fields := reflect.VisibleFields(reflect.TypeOf(userToUpdate))
	// for _, field := range fields {
	// 	fmt.Printf("Key: %s\tType: %s\n", field.Name, field.Type)
	// }

	// _, err := u.dbConn.Exec(context.Background(), `UPDATE users SET first_name = $1, last_name = $2, bio = $3, email = $4, password = $5, profile = $6 WHERE id = $7`, userToUpdate.FirstName, userToUpdate.LastName, userToUpdate.Bio, userToUpdate.Email, userToUpdate.Password, userToUpdate.Profile, userId)
	// if err != nil {
	// 	return err
	// }
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
