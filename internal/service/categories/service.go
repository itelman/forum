package categories

import (
	"database/sql"
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/service/categories/adapters"
	"github.com/itelman/forum/internal/service/categories/domain"
)

type Service interface {
	CreateCategory(input *CreateCategoryInput) error
	GetAllCategories() (*GetAllCategoriesResponse, error)
	DeleteCategory(input *DeleteCategoryInput) error
}

type service struct {
	categories domain.CategoriesRepository
}

func NewService(opts ...Option) *service {
	svc := &service{}
	for _, opt := range opts {
		opt(svc)
	}

	return svc
}

type Option func(*service)

func WithSqlite(db *sql.DB) Option {
	return func(s *service) {
		s.categories = adapters.NewCategoriesRepositorySqlite(db)
	}
}

func (s *service) CreateCategory(input *CreateCategoryInput) error {
	if err := input.validate(); err != nil {
		return err
	}

	if err := s.categories.Create(domain.CreateCategoryInput{Name: input.Name}); err != nil {
		return err
	}

	return nil
}

type GetAllCategoriesResponse struct {
	Categories []*dto.Category
}

func (s *service) GetAllCategories() (*GetAllCategoriesResponse, error) {
	categories, err := s.categories.GetAll(domain.GetAllCategoriesInput{SortedByNewest: true})
	if err != nil {
		return nil, err
	}

	return &GetAllCategoriesResponse{Categories: categories}, nil
}

func (s *service) DeleteCategory(input *DeleteCategoryInput) error {
	if _, err := s.categories.Get(domain.GetCategoryInput{ID: input.ID}); err != nil {
		return err
	}

	if err := s.categories.Delete(domain.DeleteCategoryInput{ID: input.ID}); err != nil {
		return err
	}

	return nil
}
