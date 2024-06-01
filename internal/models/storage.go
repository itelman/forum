package models

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func StorageConstructor(path string) (Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return Storage{}, err
	}

	queries := []string{`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY, 
		username TEXT NOT NULL UNIQUE, 
		password TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`, `CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY, 
		user_id INTEGER REFERENCES users(id), 
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		posted_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		likes INTEGER DEFAULT 0,
		dislikes INTEGER DEFAULT 0
	);`, `CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY,
		post_id INTEGER REFERENCES posts(id),
		user_id INTEGER REFERENCES users(id), 
		content TEXT NOT NULL,
		posted_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		likes INTEGER DEFAULT 0,
		dislikes INTEGER DEFAULT 0
	);`, `CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL UNIQUE
	);`, `CREATE TABLE IF NOT EXISTS post_categories (
		post_id INTEGER REFERENCES posts(id),
		category_id INTEGER REFERENCES categories(id)
	);`, `CREATE TABLE IF NOT EXISTS post_reactions (
		post_id INTEGER REFERENCES posts(id),
		user_id INTEGER REFERENCES users(id),
		like_or_dislike INTEGER NOT NULL
	);`, `CREATE TABLE IF NOT EXISTS comment_reactions (
		comment_id INTEGER REFERENCES comments(id),
		user_id INTEGER REFERENCES users(id),
		like_or_dislike INTEGER NOT NULL
	);`}

	for _, query := range queries {
		_, err = db.Exec(query)

		if err != nil {
			return Storage{}, err
		}
	}

	return Storage{db}, nil
}

func (s *Storage) Close() {
	s.db.Close()
}

func (s *Storage) Ping() error {
	return s.db.Ping()
}
