package tag

import (
	"database/sql"
	"fmt"
	"time"
)

// Repository struct for db connection
type Repository struct {
	Conn *sql.DB
}

// Tag struct
type Tag struct {
	ID        int64     `json:"id"`
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

// SelectTagByID retrieves a product using its id
func (repository *Repository) SelectTagByID(id int64) (*Tag, error) {
	row := repository.Conn.QueryRow("SELECT t.id, t.name, t.created_at, t.updated_at "+
		"FROM tag t WHERE t.id=(?)", id)
	var name string
	var createdAt, updatedAt time.Time
	switch err := row.Scan(&id, &name, &createdAt, &updatedAt); err {
	case sql.ErrNoRows:
		return nil, sql.ErrNoRows
	case nil:
		tag := Tag{
			ID:        id,
			Name:      name,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}
		return &tag, nil
	default:
		return nil, nil
	}
}

// retrieveAllTags stored in db
func (repository *Repository) retrieveAllTags() ([]*Tag, error) {
	rows, err := repository.Conn.Query("SELECT t.id, t.name, t.created_at, " +
		"t.updated_at FROM tag t ")

	if err != nil {
		return nil, err
	}

	var id int64
	var name string
	var createdAt, updatedAt time.Time
	var tags []*Tag
	for rows.Next() {
		err := rows.Scan(&id, &name, &createdAt, &updatedAt)
		if err != nil {
			fmt.Println(err)
		}
		tags = append(tags, &Tag{
			ID:        id,
			Name:      name,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
	}

	return tags, nil
}

// insertTag posts a new tag
func (repository *Repository) insertTag(tag *Tag) error {
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

// updateTag by ID
func (repository *Repository) updateTag(tag *Tag, id int64) error {
	stmt, err := repository.Conn.Prepare("UPDATE tag SET name=(?), description=(?), " +
		"updated_at=(?) WHERE id=(?)")
	if err != nil {
		return err
	}
	var createdAt time.Time
	row := repository.Conn.QueryRow("SELECT c.created_at FROM tag c WHERE c.id=(?)", id)
	if err := row.Scan(&createdAt); err != nil {
		return err
	}
	tag.CreatedAt = createdAt
	tag.UpdatedAt = time.Now()

	_, errExec := stmt.Exec(tag.Name, tag.UpdatedAt, id)

	if errExec != nil {
		return errExec
	}

	tag.ID = id

	return nil
}

// deleteTag by ID
func (repository *Repository) deleteTag(id int64) (int64, error) {

	res, err := repository.Conn.Exec("DELETE FROM tag WHERE id=(?)", id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
