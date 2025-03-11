package notifications

import (
	"github.com/itelman/forum/internal/dto"
	"net/http"
)

func DecodeGetAllCommentNotifications(r *http.Request) interface{} {
	return &GetAllCommentNotificationsInput{dto.GetAuthUser(r).ID}
}

func DecodeGetAllPostReactionNotifications(r *http.Request) interface{} {
	return &GetAllPostReactionNotificationsInput{dto.GetAuthUser(r).ID}
}
