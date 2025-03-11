package domain

import "database/sql"

type CommentsRepository interface {
	UpdateReactionsCount(tx *sql.Tx, input UpdateCommentReactionsCountInput) error
}

type UpdateCommentReactionsCountInput struct {
	CommentID int
}
