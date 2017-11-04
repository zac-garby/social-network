package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
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
