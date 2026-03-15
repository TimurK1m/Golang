package repository

import (
	"Practice5/internal/models"
	"database/sql"
	"fmt"
)



type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetCommonFriends(user1 int, user2 int) ([]models.User, error) {

	query := `
	SELECT u.id, u.name, u.email, u.gender, u.birth_date
	FROM user_friends uf1
	JOIN user_friends uf2 
		ON uf1.friend_id = uf2.friend_id
	JOIN users u 
		ON u.id = uf1.friend_id
	WHERE uf1.user_id = $1
	AND uf2.user_id = $2
	`

	rows, err := r.db.Query(query, user1, user2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {

		var u models.User

		err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
			&u.Gender,
			&u.BirthDate,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

func (r *Repository) GetPaginatedUsers(
	page int,
	pageSize int,
	filters map[string]string,
	orderBy string,
) (models.PaginatedResponse, error) {

	var users []models.User
	offset := (page - 1) * pageSize

	query := "SELECT id, name, email, gender, birth_date FROM users WHERE 1=1"
	args := []interface{}{}
	argID := 1

	for field, value := range filters {
		if value != "" {
			query += fmt.Sprintf(" AND %s = $%d", field, argID)
			args = append(args, value)
			argID++
		}
	}

	// default sorting
	if orderBy == "" {
		orderBy = "id"
	}

	query += " ORDER BY " + orderBy

	// pagination
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argID, argID+1)
	args = append(args, pageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return models.PaginatedResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var u models.User

		err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
			&u.Gender,
			&u.BirthDate,
		)

		if err != nil {
			return models.PaginatedResponse{}, err
		}

		users = append(users, u)
	}

	// total count
	var totalCount int
	err = r.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&totalCount)
	if err != nil {
		return models.PaginatedResponse{}, err
	}

	return models.PaginatedResponse{
		Data:       users,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}