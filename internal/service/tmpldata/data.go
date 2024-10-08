package tmpldata

import (
	"forum/internal/repository/models"
	"forum/internal/service/auth/google"
	"forum/pkg/forms"
	"time"
)

type TemplateData struct {
	TemplateName      string
	CSRFToken         string
	AuthenticatedUser *models.User
	CurrentYear       int
	Flash             string
	Form              *forms.Form
	Post              *models.Post
	Comment           *models.Comment
	Posts             []*models.Post
	Comments          []*models.Comment
	Categories        []*models.Category
	PCRelations       []string
	Error             *models.Error
	Posts_Comments    []*models.Posts_Comments
	Post_Reactions    []*models.PostReaction
	Image             *models.Image
	GoogleID          string
}

func (td *TemplateData) AddDefaultData(user *models.User, flash string) {
	td.AuthenticatedUser = user
	td.CurrentYear = time.Now().Year()
	td.Flash = flash
	td.GoogleID = google.GetClientID()
}
