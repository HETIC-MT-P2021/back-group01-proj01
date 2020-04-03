package category

import (
	"database/sql"
	"fmt"
	"time"
)

// Repository struct for db connection
type Repository struct {
	Conn *sql.DB
}

// Category struct 
type Category struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Validate : interface for JSON backend validation
func (c *Category) Validate() error {

	if c.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	return nil
}

// selectCategorytByID retrieves a product using its id
func (repository *Repository) selectCategoryByID(id int64) (*Category, error) {
	row := repository.Conn.QueryRow("SELECT c.id, c.name, c.description, "+
		"c.created_at, c.updated_at FROM category c WHERE c.id=(?)", id)
	var name, description string
	var createdAt, updatedAt time.Time
	switch err := row.Scan(&id, &name, &description, &createdAt, &updatedAt); err {
	case sql.ErrNoRows:
		return nil, sql.ErrNoRows
	case nil:
		category := Category{
			ID:          id,
			Name:        name,
			Description: description,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}
		return &category, nil
	default:
		return nil, nil
	}
}
// retrieveAllCategories stored in db
func (repository *Repository) retrieveAllCategories() ([]*Category, error) {
	rows, err := repository.Conn.Query("SELECT c.id, c.name, c.description, c.created_at, " +
		"c.updated_at FROM category c ")

	if err != nil {
		return nil, err
	}

	var id int64
	var name, description string
	var createdAt, updatedAt time.Time
	var categories []*Category
	for rows.Next() {
		err := rows.Scan(&id, &name, &description, &createdAt, &updatedAt)
		if err != nil {
			fmt.Println(err)
		}
		categories = append(categories, &Category{
			ID:          id,
			Name:        name,
			Description: description,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		})
	}

	return categories, nil
}

// insertCategory posts a new category
func (repository *Repository) insertCategory(category *Category) (*Category, error) {
	stmt, err := repository.Conn.Prepare("INSERT INTO category(name, description, created_at," +
		" updated_at) VALUES(?,?,?,?)")

	if err != nil {
		return nil, err
	}
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()
	res, errExec := stmt.Exec(category.Name, category.Description, category.CreatedAt, category.UpdatedAt)

	if errExec != nil {
		return nil, errExec
	}

	lastInsertedID, errInsert := res.LastInsertId()

	if errInsert != nil {
		return nil, errInsert
	}

	category.ID = lastInsertedID

	return category, nil
}

// updateCategory by ID
func (repository *Repository) updateCategory(category *Category, id int64) (*Category,
	error) {
	stmt, err := repository.Conn.Prepare("UPDATE category SET name=(?), description=(?)," +
		"updated_at=(?) WHERE id=(?)")

	if err != nil {
		return nil, err
	}

	category.UpdatedAt = time.Now()

	_, errExec := stmt.Exec(category.Name, category.Description, category.UpdatedAt, id)

	if errExec != nil {
		return nil, errExec
	}

	category.ID = id

	return category, nil
}

// deleteCategory by ID
func (repository *Repository) deleteCategory(id int64) (int64, error) {

	res, err := repository.Conn.Exec("DELETE FROM category WHERE id=(?)", id)
	if err != nil {
		return 0, err
	} 
	return res.RowsAffected()
}
