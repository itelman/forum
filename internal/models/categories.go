package models

type Category struct {
	ID   int
	Name string
}

func (s *Storage) GetAllCategories() ([]Category, error) {
	rows, err := s.db.Query("SELECT * FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (s *Storage) GetCategoryByID(id int) (Category, error) {
	var category Category
	err := s.db.QueryRow("SELECT id, name FROM categories WHERE id = $1", id).Scan(&category.ID, &category.Name)

	return category, err
}

func (s *Storage) InsertPostCategory(post_id, category_id int) error {
	_, err := s.db.Exec("INSERT INTO post_categories (post_id, category_id) VALUES ($1, $2)", post_id, category_id)

	return err
}
