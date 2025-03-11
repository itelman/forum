package domain

import "github.com/itelman/forum/internal/dto"

type PostReactionsRepository interface {
	GetAllNotifications(input GetAllPostReactionNotificationsInput) ([]*dto.PostReaction, error)
}

type GetAllPostReactionNotificationsInput struct {
	AuthUserID     int
	SortedByNewest bool
}
