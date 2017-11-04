package server

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Zac-Garby/social-network/project"
	"github.com/Zac-Garby/social-network/user"
	"github.com/microcosm-cc/bluemonday"
	gfm "github.com/shurcooL/github_flavored_markdown"
)

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		LoggedIn *user.User
		Projects []*project.Project
	}

	projects, err := project.GetAll(s.Database)
	if err != nil {
		handleError(err, w, r)
		return
	}

	loggedIn, err := getLoggedInUser(s.Database, r)
	if err != nil {
		handleError(err, w, r)
		return
	}

	handleRequest(w, "index", "", r.URL.String(), Data{
		Projects: projects,
		LoggedIn: loggedIn,
	})
}

func (s *Server) handleLoginPage(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   "",
		Expires: time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local),
	})

	handleRequest(w, "login", "login", r.URL.String(), nil)
}

func (s *Server) handleUserPage(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		LoggedIn *user.User
		User     *user.User
		Projects []*project.Project
	}

	_, name := filepath.Split(r.URL.String())
	u, err := user.GetUserByUsername(s.Database, name)
	if err != nil {
		handleError(err, w, r)
		return
	}

	projects, err := project.GetAllByUser(s.Database, name)
	if err != nil {
		handleError(err, w, r)
		return
	}

	loggedIn, err := getLoggedInUser(s.Database, r)
	if err != nil {
		handleError(err, w, r)
		return
	}

	handleRequest(w, "user", "user", "/user", Data{
		User:     u,
		Projects: projects,
		LoggedIn: loggedIn,
	})
}

func (s *Server) handleProjectPage(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		LoggedIn     *user.User
		Project      *project.Project
		RenderedHTML template.HTML
	}

	_, id := filepath.Split(r.URL.String())
	intID, err := strconv.Atoi(id)
	if err != nil {
		handleError(err, w, r)
		return
	}

	p, err := project.GetByID(s.Database, intID)
	if err != nil {
		handleError(err, w, r)
		return
	}

	unsafeHTML := gfm.Markdown([]byte(p.TextContent))
	safeHTML := bluemonday.UGCPolicy().SanitizeBytes(unsafeHTML)

	loggedIn, err := getLoggedInUser(s.Database, r)
	if err != nil {
		handleError(err, w, r)
		return
	}

	handleRequest(w, "project", "project", "/project", Data{
		Project:      p,
		RenderedHTML: template.HTML(safeHTML),
		LoggedIn:     loggedIn,
	})
}

func (s *Server) handleNewProject(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		LoggedIn *user.User
	}

	loggedIn, err := getLoggedInUser(s.Database, r)
	if err != nil {
		handleError(err, w, r)
		return
	}

	handleRequest(w, "new-project", "new-project", "/new-project", Data{
		LoggedIn: loggedIn,
	})
}
