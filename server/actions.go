package server

import (
	"net/http"

	"github.com/Zac-Garby/social-network/session"
	"github.com/Zac-Garby/social-network/user"
)

func (s *Server) handleLogIn(w http.ResponseWriter, r *http.Request) {
	var (
		username = r.PostFormValue("username")
		password = r.PostFormValue("password")
	)

	u, err := user.Login(s.Database, username, password)
	if err != nil {
		handleError(err, w, r)
		return
	}

	sess, err := session.NewSession(s.Database, u)
	if err != nil {
		handleError(err, w, r)
		return
	}

	session.SetSessionCookie(sess, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) handleSignUp(w http.ResponseWriter, r *http.Request) {
	var (
		username    = r.PostFormValue("username")
		displayName = r.PostFormValue("displayname")
		password    = r.PostFormValue("password")
	)

	user.CreateUser(
		s.Database,
		username,
		displayName,
		password,
		"http://zealups.com/wp-content/uploads/2017/07/golang-gopher.jpg",
	)

	s.handleLogIn(w, r)
}
