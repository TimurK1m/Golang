package postgres

import "gorm.io/gorm"

type Postgres struct {
	Conn *gorm.DB
}