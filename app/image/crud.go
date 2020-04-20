package image

import (
	"database/sql"
	"fmt"
	"image_gallery/category"
	"image_gallery/helpers"
	"time"
)

// Repository struct to store db connection
type Repository struct {
	Conn *sql.DB
}

// Image struct for handling images
type Image struct {
	ID          int64              `json:"id"`
	Name        string             `json:"name"`
	Slug        string             `json:"slug"`
	Description string             `json:"description"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	CategoryID  int64              `json:"category_id"`
	Category    *category.Category `json:"category"`
}

// Validate : interface for JSON backend validation
func (i *Image) Validate() error {

	if i.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	return nil
}

func (repository *Repository) selectImageByID(id int64) (*Image, error) {
	row := repository.Conn.QueryRow(`SELECT i.id, i.name, i.slug, i.description, i.created_at, i.updated_at, i.category_id 
	FROM image i 
	WHERE i.id=?;`, id)
	var name, slug, description string
	var createdAt, updatedAt time.Time
	var categoryID int64
	err := row.Scan(&id, &name, &slug, &description, &createdAt, &updatedAt, &categoryID)
	if err != nil {
		return nil, err
	}
	image := Image{
		ID:          id,
		Name:        name,
		Slug:        slug,
		Description: description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		CategoryID:  categoryID,
	}

	return &image, nil
}

// retrieveAllImages stored in db
func (repository *Repository) retrieveAllImages() ([]*Image, error) {
	rows, err := repository.Conn.Query(`SELECT i.id, i.name, i.slug, i.description, i.created_at, i.updated_at, i.category_id, c.name, c.description
	FROM image i 
	LEFT JOIN category c 
	ON category_id = c.id;`)
	if err != nil {
		return nil, err
	}
	var images []*Image
	var id, categoryID int64
	var name, slug, description, categoryName, categoryDescription string
	var createdAt, updatedAt time.Time

	for rows.Next() {
		var image Image
		err = rows.Scan(&id, &name, &slug, &description, &createdAt, &updatedAt, &categoryID, &categoryName, &categoryDescription)
		if err != nil {
			return nil, fmt.Errorf("could not get images : %v", err)
		}
		image = Image{
			ID:          id,
			Name:        name,
			Slug:        slug,
			Description: description,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
			CategoryID:  categoryID,
			Category: &category.Category{
				Name:        categoryName,
				Description: categoryDescription,
			},
		}
		images = append(images, &image)
	}

	return images, nil
}

// insertCategory posts a new category
func (repository *Repository) insertImage(image *Image) error {

	stmt, err := repository.Conn.Prepare("INSERT INTO image(name, slug, description, created_at," +
		" updated_at, category_id) VALUES(?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	image.CreatedAt = time.Now()
	image.UpdatedAt = time.Now()

	// Generate a slug
	for {
		slug := helpers.GenerateAlphanumericToken(10)
		exists, err := repository.slugExists(slug)
		if err != nil {
			return fmt.Errorf("could not check if slug exists: %w", err)
		}

		if !exists {
			image.Slug = slug
			break
		}
	}

	res, errExec := stmt.Exec(image.Name, image.Description, image.Slug, image.CreatedAt, image.UpdatedAt,
		image.CategoryID)
	if errExec != nil {
		return fmt.Errorf("could not exec stmt: %v", errExec)
	}

	lastInsertedID, errInsert := res.LastInsertId()

	if errInsert != nil {
		return fmt.Errorf("could not retrieve last inserted ID: %v", errInsert)
	}

	image.ID = lastInsertedID

	return nil
}

// updateCategory by ID
func (repository *Repository) updateImage(category *Image, id int64) (*Image,
	error) {
	stmt, err := repository.Conn.Prepare("UPDATE category SET name=(?), description=(?), " +
		"updated_at=(?) WHERE id=(?)")
	if err != nil {
		return nil, err
	}
	var createdAt time.Time
	row := repository.Conn.QueryRow("SELECT c.created_at FROM category c WHERE c.id=?", id)
	if err := row.Scan(&createdAt); err != nil {
		return nil, err
	}
	category.CreatedAt = createdAt
	category.UpdatedAt = time.Now()

	_, errExec := stmt.Exec(category.Name, category.Description, category.UpdatedAt, id)

	if errExec != nil {
		return nil, errExec
	}

	category.ID = id

	//TODO(athenais) fix created at

	return category, nil
}

func (repository *Repository) slugExists(slug string) (bool, error) {
	row := repository.Conn.QueryRow(`SELECT i.slug FROM image i WHERE i.slug=?;`, slug)

	err := row.Scan(&slug)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}
