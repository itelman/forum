package handler

import (
	"forum/internal/service/auth"
	"forum/internal/service/tmpldata"
	"net/http"
)

func (h *Handlers) ShowCreatedPosts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/activity/created" {
		h.NotFoundHandler(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		h.ClientErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	loggedUser := auth.AuthenticatedUser(r)

	posts, err := h.App.Repository.Posts.Created(loggedUser.ID)
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	if loggedUser != nil {
		for _, post_comments := range posts {
			reacted, err := h.App.Repository.Post_Reactions.Get(post_comments.Post.ID, loggedUser.ID)
			if err != nil {
				h.ServerErrorHandler(w, r, err)
				return
			}

			post_comments.Post.ReactedByUser = reacted
		}
	}

	err = h.App.Render(w, r, &tmpldata.TemplateData{
		TemplateName:   "activity_page.html",
		Posts_Comments: posts,
	})
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}
}

func (h *Handlers) ShowReactedPosts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/activity/reacted" {
		h.NotFoundHandler(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		h.ClientErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	loggedUser := auth.AuthenticatedUser(r)

	posts, err := h.App.Repository.Posts.Reacted(loggedUser.ID, h.App.Repository.Post_Reactions.GetReactionsByUser)
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	err = h.App.Render(w, r, &tmpldata.TemplateData{
		TemplateName:   "activity_page.html",
		Posts_Comments: posts,
	})
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}
}

func (h *Handlers) ShowCommentedPosts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/activity/commented" {
		h.NotFoundHandler(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		h.ClientErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	loggedUser := auth.AuthenticatedUser(r)

	posts, err := h.App.Repository.Posts.Commented(loggedUser.ID, h.App.Repository.Comments.GetDistinctCommentsByUser)
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	if loggedUser != nil {
		for _, post_comments := range posts {
			reacted, err := h.App.Repository.Post_Reactions.Get(post_comments.Post.ID, loggedUser.ID)
			if err != nil {
				h.ServerErrorHandler(w, r, err)
				return
			}

			comments, err := h.App.Repository.Comments.GetByUserForPost(post_comments.Post.ID, loggedUser.ID)
			if err != nil {
				h.ServerErrorHandler(w, r, err)
				return
			}

			for _, comment := range comments {
				reacted, err := h.App.Repository.Comment_Reactions.Get(comment.ID, loggedUser.ID)
				if err != nil {
					h.ServerErrorHandler(w, r, err)
					return
				}

				comment.ReactedByUser = reacted
			}

			post_comments.Comments = comments
			post_comments.Post.ReactedByUser = reacted
		}
	}

	err = h.App.Render(w, r, &tmpldata.TemplateData{
		TemplateName:   "activity_page.html",
		Posts_Comments: posts,
	})
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}
}
