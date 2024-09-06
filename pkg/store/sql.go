package store

import (
	"database/sql"
	"os"
	"path/filepath"
)

func NewSQL(name, dsn string) (*sql.DB, error) {
	dir := filepath.Dir(dsn)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(name, dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
