package server

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/microcosm-cc/bluemonday"
	gfm "github.com/shurcooL/github_flavored_markdown"

	"github.com/gorilla/mux"

	"github.com/Zac-Garby/social-network/project"
	"github.com/Zac-Garby/social-network/session"
	"github.com/Zac-Garby/social-network/user"
)

// A Server hosts the website.
type Server struct {
	Database *sql.DB
}

// Start starts the server listening on the
// given address.
func (s *Server) Start(addr string) {
	r := mux.NewRouter()

	path, err := os.Getwd()
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Pages/resources
	r.PathPrefix("/public/").Handler(http.FileServer(http.Dir(path)))
	r.HandleFunc("/", s.handleIndex)
	r.HandleFunc("/login", s.handleLoginPage)
	r.HandleFunc("/u/{user}", s.handleUserPage)
	r.HandleFunc("/p/{page}", s.handleProjectPage)

	// Form actions
	r.HandleFunc("/log-in", s.handleLogIn).Methods("POST")
	r.HandleFunc("/sign-up", s.handleSignUp).Methods("POST")

	http.Handle("/", r)

	fmt.Printf("listening on %s\n", addr)
	http.ListenAndServe(addr, nil)
}

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
