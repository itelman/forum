package post_reactions

import (
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/service/post_reactions/domain"
	"net/http"
	"strconv"
)

func DecodeCreatePostReaction(r *http.Request) (interface{}, error) {
	if err := r.ParseForm(); err != nil {
		return nil, domain.ErrPostReactionsBadRequest
	}

	postId, err := strconv.Atoi(r.PostForm.Get("post_id"))
	if err != nil {
		return nil, domain.ErrPostReactionsBadRequest
	}

	isLike, err := strconv.Atoi(r.PostForm.Get("is_like"))
	if err != nil {
		return nil, domain.ErrPostReactionsBadRequest
	}

	if !(isLike == 0 || isLike == 1) {
		return nil, domain.ErrPostReactionsBadRequest
	}

	return &CreatePostReactionInput{
		PostID: postId,
		UserID: dto.GetAuthUser(r).ID,
		IsLike: isLike,
	}, nil
}
