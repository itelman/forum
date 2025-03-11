package activity

type GetAllCreatedPostsInput struct {
	AuthUserID int
}

type GetAllReactedPostsInput struct {
	AuthUserID int
}

type GetAllCommentedPostsInput struct {
	AuthUserID int
}
