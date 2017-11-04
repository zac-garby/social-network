package user

import (
	"regexp"

	"github.com/pkg/errors"
)

var (
	// ErrInvalidUsername means that the user's username is invalid.
	ErrInvalidUsername = errors.New("user: invalid username")

	// ErrInvalidDisplayName means that the user's display name is invalid.
	ErrInvalidDisplayName = errors.New("user: invalid display name")
)

// Validate checks various fields in the user,
// and returns an error if any of them is invalid. Otherwise,
// it returns nil.
func (u *User) Validate() error {
	var (
		usernameRegex = regexp.MustCompile(`^[a-zA-Z-_.]{1,32}$`)
		displayRegex  = regexp.MustCompile(`[\\u0028-\\u00A5 &!"']{1,48}`)
	)

	if !usernameRegex.MatchString(u.Username) {
		return ErrInvalidUsername
	}

	if !displayRegex.MatchString(u.DisplayName) {
		return ErrInvalidDisplayName
	}

	return nil
}
