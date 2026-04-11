package usecase

import "Secure/internal/entity"

type (
	UserInterface interface {
		LoginUser(user *entity.LoginUserDTO) (string, error) 
		RegisterUser(user *entity.User) (*entity.User, string, error)
		GetMe(userID string) (*entity.User, error)
		PromoteUser(id string) error
	}
)