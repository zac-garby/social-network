package server

import (
	"log"
	"net/http"

	"github.com/Zac-Garby/social-network/user"
)

func handleError(err error, w http.ResponseWriter, r *http.Request) {
	log.Println(err)

	switch err {
	case ErrSessionTimedOut, user.ErrWrongPassword, user.ErrUserAlreadyExists:
		http.Redirect(w, r, "/login", http.StatusSeeOther)

	default:
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
