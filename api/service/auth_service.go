package service

import (
	"github.com/chirag1807/task-management-system/api/model/request"
	"github.com/chirag1807/task-management-system/api/model/response"
	"github.com/chirag1807/task-management-system/api/repository"
)

type AuthService interface {
	UserRegistration(user request.User) (int64, error)
	UserLogin(user request.UserCredentials) (response.User, string, error)
	RefreshToken(refreshToken string) (int64, string, error)
}

type authService struct {
	authRepository repository.AuthRepository
}

func NewAuthService(authRepository repository.AuthRepository) AuthService {
	return authService{
		authRepository: authRepository,
	}
}

func (a authService) UserRegistration(user request.User) (int64, error) {
	return a.authRepository.UserRegistration(user)
}

func (a authService) UserLogin(user request.UserCredentials) (response.User, string, error) {
	return a.authRepository.UserLogin(user)
}

func (a authService) RefreshToken(refreshToken string) (int64, string, error) {
	return a.authRepository.RefreshToken(refreshToken)
}
