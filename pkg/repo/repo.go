package repo

import "database/sql"

type Config struct {
	Connection string
}

type Connection struct {
	read  *sql.DB
	write *sql.DB
}

func NewConnection(config Config) (*Connection, error) {
	read, err := sql.Open("sqlite", config.Connection)
	if err != nil {
		return nil, err
	}
	write, err := sql.Open("sqlite", config.Connection)
	if err != nil {
		return nil, err
	}
	return &Connection{read, write}, nil
}
