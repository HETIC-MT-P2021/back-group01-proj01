package tag

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// Repository struct for db connection
type Repository struct {
	Conn *sql.DB
}

// Tag struct
type Tag struct {
	ID        int64     `json:"id,omitempty"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate : interface for JSON backend validation
func (t *Tag) Validate() error {

	if t.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if len(t.Name) > 255 {
		return fmt.Errorf("name cannot be longer than 255 characters")
	}
	return nil
}

// SelectTagBy retrieves a tag by any field (whereColumn) and any value (whereValue)
func (repository *Repository) SelectTagBy(whereColumn string, whereValue interface{}) (*Tag, error) {
	query := "SELECT t.id, t.name, t.created_at, t.updated_at FROM tag t WHERE t." + whereColumn + "=(?)"
	log.Printf("query :%v", query)
	row := repository.Conn.QueryRow(query, whereValue)
	var id int64
	var name string
	var createdAt, updatedAt time.Time
	switch err := row.Scan(&id, &name, &createdAt, &updatedAt); err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		tag := Tag{
			ID:        id,
			Name:      name,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}
		log.Printf("tag:%v", tag)
		return &tag, nil
	default:
		return nil, err
	}
}

// InsertTag posts a new tag
func (repository *Repository) InsertTag(tag *Tag) error {
	stmt, err := repository.Conn.Prepare("INSERT INTO tag(name, created_at," +
		" updated_at) VALUES(?,?,?)")

	if err != nil {
		return err
	}
	tag.CreatedAt = time.Now()
	tag.UpdatedAt = time.Now()

	res, errExec := stmt.Exec(tag.Name, tag.CreatedAt, tag.UpdatedAt)
	if errExec != nil {
		return errExec
	}

	lastInsertedID, errInsert := res.LastInsertId()

	if errInsert != nil {
		return errInsert
	}

	tag.ID = lastInsertedID

	return nil
}

// GetAllTagsByImageID gets all tags linked to an image
func (repository *Repository) GetAllTagsByImageID(id int64) ([]string, error) {

	rows, err := repository.Conn.Query("SELECT t.name "+
		"FROM tag t INNER JOIN image_tag it ON it.tag_id = t.id "+
		"INNER JOIN image i ON it.image_id = i.id WHERE i.id = (?);", id)
	if err != nil {
		return nil, err
	}

	var name string
	var tags []string

	for rows.Next() {
		err := rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, name)
	}

	return tags, nil
}
