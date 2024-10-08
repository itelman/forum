package handler

import (
	"fmt"
	"forum/internal/repository/models"
	"forum/internal/service/auth"
	"forum/internal/service/tmpldata"
	"forum/pkg/forms"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

func (h *Handlers) CreatePostForm(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/create" {
		h.NotFoundHandler(w, r)
		return
	}

	if r.Method == http.MethodPost {
		h.CreatePost(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		h.ClientErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	c, err := h.App.Repository.Categories.Latest()
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	err = h.App.Render(w, r, &tmpldata.TemplateData{
		TemplateName: "create_page.html",
		Form:         forms.New(nil),
		Categories:   c,
	})
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}
}

func (h *Handlers) CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/create" {
		h.NotFoundHandler(w, r)
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		h.ClientErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	loggedUser := auth.AuthenticatedUser(r)

	c, err := h.App.Repository.Categories.Latest()
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	err = r.ParseMultipartForm(int64(1024 * 20))
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

	if sesStore.GetSession(sessionID).CSRFToken != r.PostFormValue("csrf_token") {
		h.ClientErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "content", "categories")

	form.MinLength("title", 5)
	form.MaxLength("title", 50)
	form.MatchesPattern("title", regexp.MustCompile(`^\S.*\S$`))

	file, handler, err := r.FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		h.ServerErrorHandler(w, r, err)
		return
	}

	if !(file == nil || handler == nil) {
		defer file.Close()

		form.ImgMaxSize(handler, 1048576*20)
		form.ImgExtension(handler, ".jpg", ".jpeg", ".png", ".gif")
	}

	if !form.Valid() {
		err = h.App.Render(w, r, &tmpldata.TemplateData{
			TemplateName: "create_page.html",
			Form:         form,
			Categories:   c,
		})
		if err != nil {
			h.ServerErrorHandler(w, r, err)
			return
		}

		return
	}

	for _, idStr := range r.PostForm["categories"] {
		id, err := strconv.Atoi(idStr)
		if err != nil || idStr != strconv.Itoa(id) {
			h.ClientErrorHandler(w, r, http.StatusBadRequest)
			return
		}

		_, err = h.App.Repository.Categories.Get(id)
		if err != nil {
			if err == models.ErrNoRecord {
				h.NotFoundHandler(w, r)
			} else {
				h.ServerErrorHandler(w, r, err)
			}
			return
		}
	}

	id, err := h.App.Repository.Posts.Insert(loggedUser.ID, form.Get("title"), form.Get("content"))
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	err = h.App.Repository.Post_Category.Insert(id, r.PostForm["categories"])
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	if !(file == nil || handler == nil) {
		filePath := fmt.Sprintf("./uploads/%s", handler.Filename)

		err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
		if err != nil {
			h.ServerErrorHandler(w, r, err)
			return
		}

		f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			h.ServerErrorHandler(w, r, err)
			return
		}
		defer f.Close()
		io.Copy(f, file)

		filePath_db := fmt.Sprintf("/uploads/%s", handler.Filename)
		err = h.App.Repository.Images.Insert(id, filePath_db)
		if err != nil {
			h.ServerErrorHandler(w, r, err)
			return
		}
	}

	sesStore.PutSessionData(sessionID, "flash", "Post successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", id), http.StatusMovedPermanently)
}

func (h *Handlers) ShowPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post" {
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

	s, err := h.App.Repository.Posts.Get(id)
	if err == models.ErrNoRecord {
		h.NotFoundHandler(w, r)
		return
	} else if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	c, err := h.App.Repository.Comments.Latest(id)
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	categories, err := h.App.Repository.Post_Category.Get(id)
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}

	image, err := h.App.Repository.Images.Get(id)
	if err != nil && err != models.ErrNoRecord {
		h.ServerErrorHandler(w, r, err)
		return
	}

	loggedUser := auth.AuthenticatedUser(r)
	if loggedUser != nil {
		reacted, err := h.App.Repository.Post_Reactions.Get(id, loggedUser.ID)
		if err != nil {
			h.ServerErrorHandler(w, r, err)
			return
		}

		s.ReactedByUser = reacted

		for _, comment := range c {
			reacted, err := h.App.Repository.Comment_Reactions.Get(comment.ID, loggedUser.ID)
			if err != nil {
				h.ServerErrorHandler(w, r, err)
				return
			}

			comment.ReactedByUser = reacted
		}
	}

	err = h.App.Render(w, r, &tmpldata.TemplateData{
		TemplateName: "show_page.html",
		Post:         s,
		Comments:     c,
		PCRelations:  categories,
		Image:        image,
	})
	if err != nil {
		h.ServerErrorHandler(w, r, err)
		return
	}
}
