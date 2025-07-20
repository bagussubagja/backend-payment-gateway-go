package services

import (
	"github.com/bagussubagja/backend-payment-gateway-go/internal/models"
	repositories "github.com/bagussubagja/backend-payment-gateway-go/internal/repositories"
)

type UserService interface {
	GetUserByID(id uint) (*models.User, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}
