package tmpldata

import (
	"forum/internal/repository/models"
	"forum/pkg/forms"
	"time"
)

type TemplateData struct {
	TemplateName      string
	AuthenticatedUser *models.User
	CurrentYear       int
	Flash             string
	Form              *forms.Form
	Post              *models.Post
	Posts             []*models.Post
	Comments          []*models.Comment
	Categories        []*models.Category
	PCRelations       []string
	Error             *models.Error
}

func (td *TemplateData) AddDefaultData(user *models.User, flash string) {
	td.AuthenticatedUser = user
	td.CurrentYear = time.Now().Year()
	td.Flash = flash
}
