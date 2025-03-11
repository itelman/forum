package domain

import "database/sql"

type PostsRepository interface {
	UpdateReactionsCount(tx *sql.Tx, input UpdatePostReactionsCountInput) error
}

type UpdatePostReactionsCountInput struct {
	PostID int
}
