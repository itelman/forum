package comment_reactions

type CreateCommentReactionInput struct {
	CommentID int
	UserID    int
	IsLike    int
}
