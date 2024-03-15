package service

import (
	"github.com/chirag1807/task-management-system/api/model/request"
	"github.com/chirag1807/task-management-system/api/model/response"
	"github.com/chirag1807/task-management-system/api/repository"
)

type UserService interface {
	GetAllPublicProfileUsers() ([]response.User, error)
	GetMyDetails(userId int64) (response.User, error)
	UpdateUserProfile(userId int64, userToUpdate request.UpdateUser) error
	SendOTPToUser(userEmail string) (int8, error)
	VerifyOTP(otpFromUser request.OTP) error
	ResetUserPassword(userEmailPassword request.User) error
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return userService{
		userRepository: userRepository,
	}
}

func (u userService) GetAllPublicProfileUsers() ([]response.User, error) {
	return u.userRepository.GetAllPublicProfileUsers()
}

func (u userService) GetMyDetails(userId int64) (response.User, error) {
	return u.userRepository.GetMyDetails(userId)
}

func (u userService) UpdateUserProfile(userId int64, userToUpdate request.UpdateUser) error {
	return u.userRepository.UpdateUserProfile(userId, userToUpdate)
}

func (u userService) SendOTPToUser(userEmail string) (int8, error) {
	return u.userRepository.SendOTPToUser(userEmail)
}

func (u userService) VerifyOTP(otpFromUser request.OTP) error {
	return u.userRepository.VerifyOTP(otpFromUser)
}

func (u userService) ResetUserPassword(userEmailPassword request.User) error {
	return u.userRepository.ResetUserPassword(userEmailPassword)
}
