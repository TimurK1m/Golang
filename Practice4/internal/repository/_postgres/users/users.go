package users

import (
	"fmt"
	"time"

	"Practice4/pkg/modules"

	"Practice4/internal/repository/_postgres"
)
type Repository struct {
db *_postgres.Dialect
executionTimeout time.Duration
}
func NewUserRepository(db *_postgres.Dialect) *Repository {
return &Repository{
db: db,
executionTimeout: time.Second * 5,
}
}
func (r *Repository) GetUsers() ([]modules.User, error) {
var users []modules.User
err := r.db.DB.Select(&users, "SELECT id, name FROM users")
if err != nil {
return nil, err
}
fmt.Println(users)
return users, nil
}