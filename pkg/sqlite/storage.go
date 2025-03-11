package sqlite

import "database/sql"

func NewSqlite(dir string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dir)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
