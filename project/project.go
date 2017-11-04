package project

import (
	"database/sql"
)

// A Project holds relevant information for a
// project, with the fields from the 'projects'
// table, and a pointer to the author User.
type Project struct {
	ID          int
	Title       string
	Description string
	TextContent string
	DateCreated string

	AuthorID       int
	AuthorName     string
	AuthorUsername string
	AuthorPicture  string
}

func query(db *sql.DB, query string, args ...interface{}) ([]*Project, error) {
	rows, err := db.Query(query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var projects []*Project

	for rows.Next() {
		proj := &Project{}

		if err := rows.Scan(
			&proj.ID,
			&proj.Title,
			&proj.Description,
			&proj.TextContent,
			&proj.DateCreated,
			&proj.AuthorID,
			&proj.AuthorName,
			&proj.AuthorUsername,
			&proj.AuthorPicture,
		); err != nil {
			return projects, err
		}

		projects = append(projects, proj)
	}

	return projects, nil
}

// GetAll returns a slice containing all
// currently defined projects in the database.
func GetAll(db *sql.DB) ([]*Project, error) {
	return query(db, `
	SELECT projects.*, users.display_name, users.user_name, users.profile_picture
	FROM projects
	INNER JOIN users
	ON projects.author = users.user_id
	ORDER BY projects.date_created DESC
	`)
}

// GetAllByUser returns a slice containing
// all projects created by a specified user.
func GetAllByUser(db *sql.DB, username string) ([]*Project, error) {
	return query(db, `
	SELECT projects.*, users.display_name, users.user_name, users.profile_picture
	FROM projects
	INNER JOIN users
	ON projects.author = users.user_id
	WHERE users.user_name = ?
	ORDER BY projects.date_created DESC`, username)
}

// GetByID returns the project whose project_id
// is equal to the id parameter.
func GetByID(db *sql.DB, id int) (*Project, error) {
	proj := &Project{}

	err := db.QueryRow(`
	SELECT projects.*, users.display_name, users.user_name, users.profile_picture
	FROM projects
	INNER JOIN users
	ON projects.author = users.user_id
	WHERE project_id = ?`, id).Scan(
		&proj.ID,
		&proj.Title,
		&proj.Description,
		&proj.TextContent,
		&proj.DateCreated,
		&proj.AuthorID,
		&proj.AuthorName,
		&proj.AuthorUsername,
		&proj.AuthorPicture,
	)

	return proj, err
}
