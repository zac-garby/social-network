package session

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/Zac-Garby/social-network/user"
)

// NewSession creates a new session ID on the
// database, and returns it.
func NewSession(db *sql.DB, u *user.User) (int, error) {
	rand.Seed(int64(time.Now().Nanosecond()))
	n := rand.Int()

	stmt, err := db.Prepare(`
	INSERT INTO sessions (
		session_id,
		user_id
	) VALUES (
		?, ?
	)`)

	if err != nil {
		return 0, err
	}

	_, err = stmt.Exec(n, u.ID)
	if err != nil {
		return 0, err
	}

	return n, nil
}

// GetUser gets the user pointed to by the
// specified session ID.
func GetUser(db *sql.DB, sess int) (*user.User, error) {
	u := &user.User{}

	err := db.QueryRow(`
	SELECT users.*
	FROM sessions
	INNER JOIN users
	ON sessions.user_id = users.user_id
	WHERE sessions.session_id = ?
	`, sess).Scan(
		&u.ID,
		&u.Username,
		&u.DisplayName,
		&u.PasswordHash,
		&u.DateJoined,
		&u.ProfilePicture,
		&u.GithubUsername,
		&u.HomepageURL,
		&u.Link1URL,
		&u.Link2URL,
		&u.Link1Name,
		&u.Link2Name,
	)

	return u, err
}

// SetSessionCookie sets the client's session_id
// cookie to the specified session ID.
func SetSessionCookie(sess int, w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:  "session_id",
		Value: fmt.Sprintf("%d", sess),

		Expires: time.Now().Add(time.Second * 600000),

		MaxAge:   600000,
		HttpOnly: false,
	}

	http.SetCookie(w, cookie)
}
