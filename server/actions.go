package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Zac-Garby/social-network/project"
	"github.com/Zac-Garby/social-network/session"
	"github.com/Zac-Garby/social-network/user"
)

// ErrUsernameInUse is returned when a user's username cannot be changed, because
// someone else already has the username they want to change to.
var ErrUsernameInUse = errors.New("change username: username already in use")

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

func (s *Server) handleAddProject(w http.ResponseWriter, r *http.Request) {
	var (
		title       = r.PostFormValue("title")
		description = r.PostFormValue("description")
		content     = r.PostFormValue("content")
	)

	if len(title) < 1 || len(title) > 75 ||
		len(description) < 1 || len(description) > 100 ||
		len(content) < 1 || len(content) > 15000 {

		http.Redirect(w, r, "/new", http.StatusSeeOther)
		return
	}

	u, err := getLoggedInUser(s.Database, r)
	if err != nil {
		handleError(err, w, r)
	}

	proj, err := project.CreateProject(
		s.Database,
		title,
		description,
		content,
		u.ID,
	)

	if err != nil {
		handleError(err, w, r)
	}

	http.Redirect(w, r, fmt.Sprintf("/p/%d", proj.ID), http.StatusSeeOther)
}

func (s *Server) handleEditProfile(w http.ResponseWriter, r *http.Request) {
	var (
		username    = r.PostFormValue("username")
		displayname = r.PostFormValue("displayname")
		profilePic  = r.PostFormValue("profilepicture")
		ghUsername  = r.PostFormValue("gh-user")
		homepage    = r.PostFormValue("homepage")
		url1name    = r.PostFormValue("linkname-1")
		url1        = r.PostFormValue("url-1")
		url2name    = r.PostFormValue("linkname-2")
		url2        = r.PostFormValue("url-2")
	)

	current, err := getLoggedInUser(s.Database, r)
	if err != nil {
		handleError(err, w, r)
		return
	}

	either := func(a, b string) string {
		if len(a) < 1 {
			return b
		}

		return a
	}

	username = either(username, current.Username)
	displayname = either(displayname, current.DisplayName)
	profilePic = either(profilePic, current.ProfilePicture)
	ghUsername = either(ghUsername, current.GithubUsername)
	homepage = either(homepage, current.HomepageURL)
	url1name = either(url1name, current.Link1Name)
	url1 = either(url1, current.Link1URL)
	url2name = either(url2name, current.Link2Name)
	url2 = either(url2, current.Link2URL)

	old, err := user.GetUserByUsername(s.Database, username)

	// User already exists with new username
	if err == nil && old.ID != current.ID {
		handleError(ErrUsernameInUse, w, r)
		return
	}

	newUser := &user.User{
		ID: current.ID,

		Username:       username,
		DisplayName:    displayname,
		ProfilePicture: profilePic,
		GithubUsername: ghUsername,
		HomepageURL:    homepage,
		Link1Name:      url1name,
		Link1URL:       url1,
		Link2Name:      url2name,
		Link2URL:       url2,
	}

	if err := user.Update(s.Database, newUser); err != nil {
		handleError(err, w, r)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/u/%s", newUser.Username), http.StatusSeeOther)
}
