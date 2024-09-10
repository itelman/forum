package handler

import (
	"fmt"
	"forum/internal/repository/models"
	"forum/internal/service/auth"
	"forum/internal/service/tmpldata"
	"forum/pkg/forms"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

func (h *Handlers) DeletePost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/delete" {
		h.NotFoundHandler(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		h.ClientErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	idQuery := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idQuery)
	if err != nil || idQuery != strconv.Itoa(id) {
		h.ClientErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	post, err := h.App.Repository.Posts.Get(id)
	if err == models.ErrNoRecord {
		h.NotFoundHandler(w, r)
		return
	} else if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	if auth.AuthenticatedUser(r).ID != post.UserID {
		h.ClientErrorHandler(w, r, http.StatusForbidden)
		return
	}

	err = h.App.Repository.Posts.Delete(id)
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	sesStore := h.App.SessionStore
	sessionID, err := sesStore.GetSessionIDFromRequest(w, r)
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}
	sesStore.PutSessionData(sessionID, "flash", "Post successfully removed!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handlers) EditPostForm(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/edit" {
		h.NotFoundHandler(w, r)
		return
	}

	if r.Method == http.MethodPost {
		h.EditPost(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		h.ClientErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	idQuery := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idQuery)
	if err != nil || idQuery != strconv.Itoa(id) {
		h.ClientErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	post, err := h.App.Repository.Posts.Get(id)
	if err == models.ErrNoRecord {
		h.NotFoundHandler(w, r)
		return
	} else if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	if auth.AuthenticatedUser(r).ID != post.UserID {
		h.ClientErrorHandler(w, r, http.StatusForbidden)
		return
	}

	getForm := make(url.Values)
	getForm.Set("title", post.Title)
	getForm.Set("content", post.Content)

	err = h.App.Render(w, r, &tmpldata.TemplateData{
		TemplateName: "editpost_page.html",
		Post:         post,
		Form:         forms.New(getForm),
	})
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}
}

func (h *Handlers) EditPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/edit" {
		h.NotFoundHandler(w, r)
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		h.ClientErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	idQuery := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idQuery)
	if err != nil || idQuery != strconv.Itoa(id) {
		h.ClientErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	post, err := h.App.Repository.Posts.Get(id)
	if err == models.ErrNoRecord {
		h.NotFoundHandler(w, r)
		return
	} else if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	if auth.AuthenticatedUser(r).ID != post.UserID {
		h.ClientErrorHandler(w, r, http.StatusForbidden)
		return
	}

	err = r.ParseForm()
	if err != nil {
		h.ClientErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "content")

	form.MinLength("title", 5)
	form.MaxLength("title", 50)
	form.MatchesPattern("title", regexp.MustCompile(`^\S.*\S$`))

	title := form.Get("title")
	content := form.Get("content")

	if !form.Valid() {
		err = h.App.Render(w, r, &tmpldata.TemplateData{
			Post:         post,
			TemplateName: "editpost_page.html",
			Form:         form,
		})
		if err != nil {
			h.ServerErrorHandler(w, r, err)
			return
		}

		return
	}

	if title == post.Title && content == post.Content {
		form.Errors.Add("generic", "Please type something new or cancel edit.")
		err = h.App.Render(w, r, &tmpldata.TemplateData{
			Post:         post,
			TemplateName: "editpost_page.html",
			Form:         form,
		})
		if err != nil {
			h.ServerErrorHandler(w, r, err)
			return
		}

		return
	}

	err = h.App.Repository.Posts.Update(id, title, content)
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	sesStore := h.App.SessionStore
	sessionID, err := sesStore.GetSessionIDFromRequest(w, r)
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}
	sesStore.PutSessionData(sessionID, "flash", "Post successfully edited!")
	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", id), http.StatusSeeOther)
}

func (h *Handlers) DeleteComment(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/comment/delete" {
		h.NotFoundHandler(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		h.ClientErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	idQuery := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idQuery)
	if err != nil || idQuery != strconv.Itoa(id) {
		h.ClientErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	comment, err := h.App.Repository.Comments.Get(id)
	if err == models.ErrNoRecord {
		h.NotFoundHandler(w, r)
		return
	} else if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	if auth.AuthenticatedUser(r).ID != comment.UserID {
		h.ClientErrorHandler(w, r, http.StatusForbidden)
		return
	}

	err = h.App.Repository.Comments.Delete(id)
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	sesStore := h.App.SessionStore
	sessionID, err := sesStore.GetSessionIDFromRequest(w, r)
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}
	sesStore.PutSessionData(sessionID, "flash", "Comment successfully removed!")
	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", comment.PostID), http.StatusSeeOther)
}

func (h *Handlers) EditCommentForm(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/comment/edit" {
		h.NotFoundHandler(w, r)
		return
	}

	if r.Method == http.MethodPost {
		h.EditComment(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		h.ClientErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	idQuery := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idQuery)
	if err != nil || idQuery != strconv.Itoa(id) {
		h.ClientErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	comment, err := h.App.Repository.Comments.Get(id)
	if err == models.ErrNoRecord {
		h.NotFoundHandler(w, r)
		return
	} else if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	if auth.AuthenticatedUser(r).ID != comment.UserID {
		h.ClientErrorHandler(w, r, http.StatusForbidden)
		return
	}

	getForm := make(url.Values)
	getForm.Set("content", comment.Content)

	err = h.App.Render(w, r, &tmpldata.TemplateData{
		TemplateName: "editcomment_page.html",
		Comment:      comment,
		Form:         forms.New(getForm),
	})
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}
}

func (h *Handlers) EditComment(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/comment/edit" {
		h.NotFoundHandler(w, r)
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		h.ClientErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	idQuery := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idQuery)
	if err != nil || idQuery != strconv.Itoa(id) {
		h.ClientErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	comment, err := h.App.Repository.Comments.Get(id)
	if err == models.ErrNoRecord {
		h.NotFoundHandler(w, r)
		return
	} else if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	if auth.AuthenticatedUser(r).ID != comment.UserID {
		h.ClientErrorHandler(w, r, http.StatusForbidden)
		return
	}

	err = r.ParseForm()
	if err != nil {
		h.ClientErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("content")

	content := form.Get("content")

	if !form.Valid() {
		err = h.App.Render(w, r, &tmpldata.TemplateData{
			TemplateName: "editcomment_page.html",
			Comment:      comment,
			Form:         form,
		})
		if err != nil {
			h.ServerErrorHandler(w, r, err)
			return
		}
		return
	}

	if content == comment.Content {
		form.Errors.Add("generic", "Please type something new or cancel edit.")
		err = h.App.Render(w, r, &tmpldata.TemplateData{
			TemplateName: "editcomment_page.html",
			Comment:      comment,
			Form:         form,
		})
		if err != nil {
			h.ServerErrorHandler(w, r, err)
			return
		}

		return
	}

	err = h.App.Repository.Comments.Update(id, content)
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	sesStore := h.App.SessionStore
	sessionID, err := sesStore.GetSessionIDFromRequest(w, r)
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}
	sesStore.PutSessionData(sessionID, "flash", "Comment successfully edited!")
	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", comment.PostID), http.StatusSeeOther)
}
