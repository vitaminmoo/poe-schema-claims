package api

import "database/sql"

// StrictServerInterface being defined like this causes unimplemented handler methods to be a compilation error.
var _ StrictServerInterface = (*server)(nil)

type server struct {
	db *sql.DB
}

func NewServer(db *sql.DB) server {
	return server{
		db: db,
	}
}
