package user

import (
	"database/sql"

	"github.com/pkg/errors"
)

// ErrUserAlreadyExists signifies that a new user cannot
// be created, because another user of the same username
// already exists in the database.
var ErrUserAlreadyExists = errors.New("sign-up: user already exists")

// CreateUser makes a new user with the given parameters,
// and adds the new user to the database. If a user
// already exists with the same username, it won't
// create a new user, and will return ErrUserAlreadyExists.
func CreateUser(db *sql.DB, username, displayname, password, profilePicture string) (*User, error) {
	_, err := GetUserByUsername(db, username)

	if err == nil {
		return nil, ErrUserAlreadyExists
	}

	if err != sql.ErrNoRows {
		return nil, err
	}

	hashed, err := hash(password)
	if err != nil {
		return nil, err
	}

	stmt, err := db.Prepare(`INSERT INTO users (
		user_id,
		user_name,
		display_name,
		password_hash,
		date_joined,
		profile_picture
	) VALUES (
		NULL, ?, ?, ?, NOW(), ?
	)`)

	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(username, displayname, hashed, profilePicture)
	if err != nil {
		return nil, err
	}

	return GetUserByUsername(db, username)
}
