package entities

import (
	"time"
)

// User represents a user entity
type User struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

// Validate validates user data
func (u *User) Validate() error {
	if u.Name == "" {
		return ErrUserNameRequired
	}
	if u.Email == "" {
		return ErrUserEmailRequired
	}
	return nil
}

// IsValid returns true if user data is valid
func (u *User) IsValid() bool {
	return u.Validate() == nil
}

// UpdateName updates the user's name
func (u *User) UpdateName(name string) error {
	if name == "" {
		return ErrUserNameRequired
	}
	u.Name = name
	u.Updated = time.Now()
	return nil
}

// UpdateEmail updates the user's email
func (u *User) UpdateEmail(email string) error {
	if email == "" {
		return ErrUserEmailRequired
	}
	u.Email = email
	u.Updated = time.Now()
	return nil
}