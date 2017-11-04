package user

import (
	"crypto/sha512"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

// ErrWrongPassword signifies that a login attempt
// failed, because the user entered the wrong password.
var ErrWrongPassword = errors.New("login: wrong password")

// A User contains all the information stored
// for a single user.
type User struct {
	ID             int
	Username       string
	DisplayName    string
	PasswordHash   string
	DateJoined     string
	ProfilePicture string
}

// GetUser queries the database for a user whose
// id matches the id argument.
func GetUser(db *sql.DB, id int) (*User, error) {
	user := &User{}

	err := db.QueryRow("SELECT * FROM users WHERE user_id = ?", id).Scan(
		&user.ID,
		&user.Username,
		&user.DisplayName,
		&user.PasswordHash,
		&user.DateJoined,
		&user.ProfilePicture,
	)

	return user, err
}

// GetUserByUsername queries the database for a user whose
// username matches the username argument.
func GetUserByUsername(db *sql.DB, username string) (*User, error) {
	user := &User{}

	err := db.QueryRow("SELECT * FROM users WHERE user_name = ?", username).Scan(
		&user.ID,
		&user.Username,
		&user.DisplayName,
		&user.PasswordHash,
		&user.DateJoined,
		&user.ProfilePicture,
	)

	return user, err
}

// GetAllUsers returns a slice containing all users in the
// database, ordered by display name.
func GetAllUsers(db *sql.DB) ([]*User, error) {
	rows, err := db.Query("SELECT * FROM users ORDER BY users.display_name")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*User

	for rows.Next() {
		user := &User{}

		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.DisplayName,
			&user.PasswordHash,
			&user.DateJoined,
			&user.ProfilePicture,
		); err != nil {
			return users, err
		}

		users = append(users, user)
	}

	return users, nil
}

// Login finds the user with the specified username, and
// returns the user if the hash of the specified password
// equals the user's password, otherwise returns an error.
func Login(db *sql.DB, username, password string) (*User, error) {
	user, err := GetUserByUsername(db, username)
	if err != nil {
		return nil, err
	}

	hashed, err := hash(password)
	if err != nil {
		return nil, err
	}

	if hashed == user.PasswordHash {
		return user, nil
	}

	return nil, ErrWrongPassword
}

// Update updates the user name, display name, and profile
// picture fields of the specified user.
func Update(db *sql.DB, u *User) error {
	stmt, err := db.Prepare(`
	UPDATE users
	SET user_name = ?,
		display_name = ?,
		profile_picture = ?
	WHERE users.user_id = ?`)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(u.Username, u.DisplayName, u.ProfilePicture, u.ID)
	if err != nil {
		return err
	}

	return nil
}

func hash(str string) (string, error) {
	hasher := sha512.New()
	_, err := hasher.Write([]byte(str))
	if err != nil {
		return "", err
	}

	hashed := hasher.Sum(nil)

	return fmt.Sprintf("%x", hashed), nil
}
