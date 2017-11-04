package server

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/Zac-Garby/social-network/session"
	"github.com/Zac-Garby/social-network/user"
)

// ErrSessionTimedOut is returned from getLoggedInUser
// when no session cookie is found, because it's assumed
// that that's the reason.
var ErrSessionTimedOut = errors.New("user: session timed out")

func getLoggedInUser(db *sql.DB, r *http.Request) (*user.User, error) {
	cookies := r.Cookies()

	for _, cookie := range cookies {
		if cookie.Name == "session_id" {
			sess, err := strconv.Atoi(cookie.Value)
			if err != nil {
				return nil, err
			}

			return session.GetUser(db, sess)
		}
	}

	return nil, ErrSessionTimedOut
}
