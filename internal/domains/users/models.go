package users

import (
	"errors"
	"time"
)

const (
	RoleUser = "USER"
)

var (
	ErrBlankPassword       = errors.New("password is blank")
	ErrPasswordsDoNotMatch = errors.New("passwords do not match")
	ErrUsernameRegistered  = errors.New("username is already registered")
	ErrUserNotFound        = errors.New("username is not registered")
	ErrPasswordInvalid     = errors.New("wrong password")
)

type User struct {
	ID           string    `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Roles        []string  `db:"roles" json:"roles"`
	Accounts     []Account `db:"accounts" json:"accounts"`
}

type Account struct {
	ID       string `db:"id" json:"id"`
	UserID   string `db:"user_id" json:"-"`
	Service  string `db:"service" json:"service"`
	Username string `db:"username" json:"username"`

	RefreshStartedAt   *time.Time `db:"refresh_started_at" json:"refresh_started_at"`
	RefreshActiveAt    *time.Time `db:"refresh_active_at" json:"refresh_active_at"`
	RefreshCompletedAt *time.Time `db:"refresh_completed_at" json:"refresh_completed_at"`
	RefreshStatus      string     `db:"refresh_status" json:"refresh_status"`
}

type NewUser struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type NewAccount struct {
	Service  string `json:"service"`
	Username string `json:"username"`
}
