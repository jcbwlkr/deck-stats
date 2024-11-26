package users

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type User struct {
	ID           string    `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	Password     string    `db:"password" json:"-"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Roles        []string  `db:"roles" json:"roles"`
	Accounts     []Account `db:"accounts" json:"accounts"`
}

type Account struct {
	ID       string `db:"id" json:"id"`
	Service  string `db:"service" json:"service"`
	Token    string `db:"token" json:"-"`
	Username string `db:"username" json:"username"`

	RefreshStartedAt   *time.Time `db:"refresh_started_at" json:"refresh_started_at"`
	RefreshActiveAt    *time.Time `db:"refresh_active_at" json:"refresh_active_at"`
	RefreshCompletedAt *time.Time `db:"refresh_completed_at" json:"refresh_completed_at"`
	RefreshStatus      string     `db:"refresh_status" json:"refresh_status"`
}

type Service struct {
	db *sqlx.DB
}

func NewService(db *sqlx.DB) *Service {
	return &Service{db: db}
}

//func (s *Service) RegisterUser(ctx context.Context,

func (s *Service) UpdateAccount(ctx context.Context, account Account) error {
	const q = `
	UPDATE user_accounts SET
		service = :service
		username = :username
		token = :token
		refresh_started_at = :refresh_started_at
		refresh_active_at = :refresh_active_at
		refresh_completed_at = :refresh_completed_at
		refresh_status = :refresh_status
	WHERE id = :id`

	_, err := s.db.NamedExecContext(ctx, q, account)
	return err
}
