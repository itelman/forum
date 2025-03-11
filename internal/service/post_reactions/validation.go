package post_reactions

type CreatePostReactionInput struct {
	PostID int
	UserID int
	IsLike int
}
