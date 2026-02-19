package repository

import (
	"Practice4/internal/repository/_postgres"
	"Practice4/internal/repository/_postgres/users"
	"Practice4/pkg/modules"
)
type UserRepository interface {
	GetUsers() ([]modules.User, error)
	GetUserByID(id int64) (*modules.User, error)
	CreateUser(user *modules.User) (int64, error)
	UpdateUser(user *modules.User) error
	DeleteUser(id int64) (int64, error)
}

type Repositories struct {
	UserRepository
}
func NewRepositories(db *_postgres.Dialect) *Repositories {
	return &Repositories{
		UserRepository: users.NewUserRepository(db),
	}
}