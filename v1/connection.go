package v1

import (
	"github.com/jmoiron/sqlx"
)

type InDB struct {
	DB *sqlx.DB
}