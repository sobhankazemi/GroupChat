package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const maxOpenDBConnection = 10
const maxIdleDBConnection = 5
const maxDBLifeTime = 5 * time.Minute

func NewConnection(connectionString string) (*sql.DB, error) {
	d, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, err
	}

	if err = d.Ping(); err != nil {
		return nil, err
	}
	d.SetConnMaxLifetime(maxDBLifeTime)
	d.SetConnMaxIdleTime(maxIdleDBConnection)
	d.SetMaxOpenConns(maxOpenDBConnection)

	return d, nil
}
