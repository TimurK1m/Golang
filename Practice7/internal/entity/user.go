package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
type User struct{
	ID	   uuid.UUID `json:"id"`
	Username string
	Email string
	Password string
	Role string
	Verified bool
}
