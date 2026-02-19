package users

import (
	"database/sql"
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
	err := r.db.DB.Select(&users, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) GetUserByID(id int64) (*modules.User, error) {
	var user modules.User

	err := r.db.DB.Get(&user, "SELECT * FROM users WHERE id=$1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, err
	}

	return &user, nil
}

func (r *Repository) CreateUser(user *modules.User) (int64, error) {
	query := `
	INSERT INTO users (name, email, age, address)
	VALUES ($1, $2, $3, $4)
	RETURNING id
	`

	var id int64
	err := r.db.DB.QueryRow(
		query,
		user.Name,
		user.Email,
		user.Age,
		user.Address,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) UpdateUser(user *modules.User) error {
	query := `
	UPDATE users
	SET name=$1, email=$2, age=$3, address=$4
	WHERE id=$5
	`

	result, err := r.db.DB.Exec(
		query,
		user.Name,
		user.Email,
		user.Age,
		user.Address,
		user.ID,
	)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *Repository) DeleteUser(id int64) (int64, error) {
	result, err := r.db.DB.Exec(
		"DELETE FROM users WHERE id=$1",
		id,
	)
	if err != nil {
		return 0, err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return 0, fmt.Errorf("user not found")
	}

	return rows, nil
}
