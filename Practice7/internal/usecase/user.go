package usecase

import (
	"Secure/internal/entity"
	"Secure/internal/usecase/repo"
	"Secure/utils"
	"fmt"

	"github.com/google/uuid"
)

type UserUseCase struct {
	repo *repo.UserRepo
}

type UserRepo struct {
	repo *repo.UserRepo
}

func NewUserUseCase(r *repo.UserRepo) *UserUseCase {
	return &UserUseCase{
		repo: r,
	}
}

func (u *UserUseCase) PromoteUser(id string) error {
	return u.repo.PromoteToAdmin(id)
}

func (u *UserUseCase) RegisterUser(user *entity.User) (*entity.User, string, error) {
	user, err := u.repo.RegisterUser(user)
	if err != nil {
		return nil, "", fmt.Errorf("Register user: %w", err)
	}

	sessionID := uuid.New().String()
	return user, sessionID, nil
}

func (u *UserUseCase) LoginUser(user *entity.LoginUserDTO) (string, error){
	userFromRepo, err := u.repo.LoginUser(user)
	if err != nil {
		return "", fmt.Errorf("User From Repo: %w", err)
	}
	if !utils.CheckPassword(userFromRepo.Password, user.Password) {
		return "", fmt.Errorf("Check Password: %w", err)
	}

	token, err := utils.GenerateJWT(userFromRepo.ID, userFromRepo.Role)
	if err != nil {
		return "", fmt.Errorf("Generate JWT: %w", err)
	}
	return token, nil
}

func (u *UserUseCase) GetMe(userID string) (*entity.User, error) {
	return u.repo.GetByID(userID)
}
