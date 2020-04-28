package image

import (
	"database/sql"
	"fmt"
	"image_gallery/category"
	"image_gallery/helpers"
	"image_gallery/tag"
	"strings"
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
	Type        string             `json:"type,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	CategoryID  int64              `json:"category_id"`
	Category    *category.Category `json:"category"`
	Tags        []*tag.Tag         `json:"tags"`
}

// Validate : interface for JSON backend validation
func (i *Image) Validate() error {

	if i.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	return nil
}

func (repository *Repository) selectImageByID(id int64) (*Image, error) {
	row := repository.Conn.QueryRow(`SELECT i.id, i.name, i.slug, i.description, i.type, 
	i.created_at, i.updated_at, i.category_id FROM image i WHERE i.id=?;`, id)
	var name, slug, description, typeExt string
	var createdAt, updatedAt time.Time
	var categoryID int64
	switch err := row.Scan(&id, &name, &slug, &description, &typeExt, &createdAt, &updatedAt, &categoryID); err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		image := Image{
			ID:          id,
			Name:        name,
			Slug:        slug,
			Description: description,
			Type:        typeExt,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
			CategoryID:  categoryID,
		}
		return &image, nil
	default:
		return nil, err
	}
}

type filterName string

const filterByDateOfUpdate filterName = "updated_at"
const filterByTag filterName = "tag"
const filterByCategory filterName = "category"

// retrieveAllImages stored in db
func (repository *Repository) retrieveAllImages(filters map[filterName]interface{}) ([]*Image, error) {

	queryFilters := make([]string, 0)
	queryArgs := make([]interface{}, 0)
	queryJoins := make([]string, 0)
	queryOrders := make([]string, 0)
	scan := make([]interface{}, 0)

	var id, categoryID int64
	var name, slug, description, typeExt, categoryName, categoryDescription, tagName string
	var createdAt, updatedAt time.Time

	scan = append(scan, &id, &name, &slug, &description, &typeExt, &createdAt, &updatedAt, &categoryID)

	queryFields := []string{
		"i.id", "i.name", "i.slug", "i.description", "i.type", "i.created_at", "i.updated_at", "i.category_id",
	}

	// Filtering lessons by date
	if v, ok := filters[filterByDateOfUpdate]; ok {
		if vv, ok := v.(string); ok {
			switch vv {
			case "asc":
				queryOrders = append(queryOrders, "updated_at ASC")
			case "desc":
				queryOrders = append(queryOrders, "updated_at DESC")
			}
		}
	}

	if v, ok := filters[filterByCategory]; ok {
		if vv, ok := v.(int64); ok {
			queryFilters = append(queryFilters, "i.category_id = ?")
			queryArgs = append(queryArgs, vv)
			queryJoins = append(queryJoins, "INNER JOIN category c ON c.id = i.category_id")
			queryFields = append(queryFields, "c.name AS category_name", "c.description AS"+
				" category_description")
			scan = append(scan, &categoryName, &categoryDescription)
		}
	}

	if v, ok := filters[filterByTag]; ok {
		if vv, ok := v.(int64); ok {
			queryFilters = append(queryFilters, "t.id = ?")
			queryArgs = append(queryArgs, vv)
			queryJoins = append(queryJoins, "INNER JOIN image_tag it ON it.image_id = i.id"+
				" INNER JOIN tag t ON t.id = it.tag_id")
			queryFields = append(queryFields, "t.name AS tag_name")
			scan = append(scan, &tagName)
		}
	}

	query := fmt.Sprintf("SELECT %s FROM image i %s",
		strings.Join(queryFields, ", "), strings.Join(queryJoins, "\n"))

	if len(queryFilters) > 0 {
		query += fmt.Sprintf("\nWHERE %s", strings.Join(queryFilters, "\nAND "))
	}

	if len(queryOrders) > 0 {
		query += fmt.Sprintf("\nORDER BY %s", strings.Join(queryOrders, ", "))
	}

	fmt.Printf("QUERY :\n %s", query)

	rows, err := repository.Conn.Query(query, queryArgs...)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve images: %v", err)
	}

	var images []*Image

	for rows.Next() {

		var image Image
		err = rows.Scan(scan...)
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
		}

		if categoryName != "" {
			image.Category = &category.Category{
				Name:        categoryName,
				Description: categoryDescription,
			}
		}

		if tagName != "" {
			image.Tags = []*tag.Tag{{Name: tagName}}
		}

		images = append(images, &image)
	}

	return images, nil
}

// insertCategory posts a new image
func (repository *Repository) insertImage(image *Image) error {

	stmt, err := repository.Conn.Prepare("INSERT INTO image(name, slug, description, type, created_at," +
		" updated_at, category_id) VALUES(?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}

	image.Type = ""
	image.CreatedAt = time.Now()
	image.UpdatedAt = time.Now()

	// Generate a slug
	for {
		slug := helpers.GenerateAlphanumericToken(10)
		id, err := repository.checkIfRowExists("image", "slug", slug)
		if err != nil {
			return fmt.Errorf("could not check if slug exists: %w", err)
		}

		if id == 0 {
			image.Slug = slug
			break
		}
	}

	res, errExec := stmt.Exec(image.Name, image.Slug, image.Description, image.Type, image.CreatedAt, image.UpdatedAt,
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

// updateImage by ID
func (repository *Repository) updateImage(image *Image, id int64) error {
	stmt, err := repository.Conn.Prepare("UPDATE image SET name=(?), description=(?), type=(?)," +
		"updated_at=(?) WHERE id=(?)")
	if err != nil {
		return err
	}
	var slug string
	var createdAt time.Time
	row := repository.Conn.QueryRow(`SELECT i.slug, i.created_at FROM image i WHERE i.id=?`, id)
	if err := row.Scan(&slug, &createdAt); err != nil {
		return err
	}

	image.CreatedAt = createdAt
	image.Slug = slug
	image.UpdatedAt = time.Now()

	_, errExec := stmt.Exec(image.Name, image.Description, image.Type, image.UpdatedAt, id)

	if errExec != nil {
		return errExec
	}

	image.ID = id

	return nil
}

// deleteCategory by ID
func (repository *Repository) deleteImage(id int64) (int64, error) {

	res, err := repository.Conn.Exec("DELETE FROM image WHERE id=(?)", id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (repository *Repository) checkIfRowExists(tableName string, WhereColumn string, whereValue interface{}) (int64, error) {
	row := repository.Conn.QueryRow("SELECT id FROM "+tableName+" WHERE "+WhereColumn+"=(?)", whereValue)

	var id int64

	switch err := row.Scan(&whereValue); err {
	case sql.ErrNoRows:
		return 0, nil
	case nil:
		return id, nil
	default:
		return 0, err
	}
}

// add tag id and image id to Many To Many Table
func (repository *Repository) linkTagToImage(imageID int64, tagID int64) error {

	stmt, err := repository.Conn.Prepare("INSERT INTO image_tag(image_id, tag_id)" +
		"VALUES(?,?)")
	if err != nil {
		return err
	}
	res, errExec := stmt.Exec(imageID, tagID)
	if errExec != nil {
		return errExec
	}

	_, errInsert := res.LastInsertId()
	if errInsert != nil {
		return errInsert
	}

	return nil
}
