package postgresrepo

import (
	"database/sql"
	"github.com/sobhankazemi/GroupChat/History/dbrepo"
)

type Repository struct {
	db *sql.DB
}

func NewPostgreRepo(_db *sql.DB) dbrepo.Repository {
	return &Repository{
		db: _db,
	}
}
