package project

import "database/sql"

// CreateProject adds a new project to the database
// with the given parameters.
func CreateProject(db *sql.DB, title, description, content string, author int) (*Project, error) {
	stmt, err := db.Prepare(`INSERT INTO projects (
		project_id,
		title,
		description,
		content,
		date_created,
		author
	) VALUES (
		NULL, ?, ?, ?, NOW(), ?
	)`)

	if err != nil {
		return nil, err
	}

	res, err := stmt.Exec(title, description, content, author)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetByID(db, int(id))
}
