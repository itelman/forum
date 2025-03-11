package domain

import "github.com/itelman/forum/internal/dto"

type CommentsRepository interface {
	GetAllNotifications(input GetAllCommentNotificationsInput) ([]*dto.Comment, error)
}

type GetAllCommentNotificationsInput struct {
	AuthUserID     int
	SortedByNewest bool
}
