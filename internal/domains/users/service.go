package users

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db *sqlx.DB
}

func NewService(db *sqlx.DB) *Service {
	return &Service{db: db}
}

func (s *Service) RegisterUser(ctx context.Context, nu NewUser) (*User, error) {

	// Sanitize and validate
	nu.Password = strings.TrimSpace(nu.Password)
	nu.PasswordConfirm = strings.TrimSpace(nu.PasswordConfirm)

	if nu.Password == "" {
		return nil, ErrBlankPassword
	}

	if nu.Password != nu.PasswordConfirm {
		return nil, ErrPasswordsDoNotMatch
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := User{
		ID:       uuid.New().String(),
		Username: nu.Username,
		Roles:    []string{RoleUser},
		Accounts: []Account{},
	}

	const q = `
		INSERT into users (id, username, password_hash, roles)
		VALUES ($1, $2, $3, $4)`

	_, err = s.db.ExecContext(ctx, q, user.ID, user.Username, string(hash), pq.StringArray(user.Roles))
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return nil, ErrUsernameRegistered
			}
		}
		return nil, err
	}

	return &user, nil
}

func (s *Service) Login(ctx context.Context, u LoginUser) (*User, error) {
	user, err := s.GetUserByUsername(ctx, u.Username)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(u.Password),
	); err != nil {
		return nil, ErrPasswordInvalid
	}

	return user, nil
}

func (s *Service) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	const q = `
	SELECT id, username, password_hash, roles
	FROM users
	WHERE username = $1`

	row := s.db.QueryRowContext(ctx, q, username)

	u, roles := User{}, pq.StringArray{}

	if err := row.Scan(
		&u.ID,
		&u.Username,
		&u.PasswordHash,
		&roles,
	); err != nil {
		return nil, err
	}
	u.Roles = roles

	return &u, nil
}

func (s *Service) GetAccount(ctx context.Context, id, userID string) (Account, error) {
	const q = `
	SELECT
		id, service, username
	FROM user_accounts
	WHERE id = $1
		AND user_id = $2`

	// TODO(jlw) this should include the refresh columns

	var a Account
	err := s.db.GetContext(ctx, &a, q, id, userID)
	return a, err
}

func (s *Service) CreateAccount(ctx context.Context, userID string, na NewAccount) (Account, error) {
	const q = `
	INSERT INTO user_accounts
	(id, user_id, service, username)
	VALUES
	(:id, :user_id, :service, :username)`

	a := Account{
		ID:       uuid.New().String(),
		UserID:   userID,
		Service:  na.Service,
		Username: na.Username,
	}

	_, err := s.db.NamedExecContext(ctx, q, a)
	return a, err
}

func (s *Service) UpdateAccount(ctx context.Context, account Account) error {
	const q = `
	UPDATE user_accounts SET
		service = :service,
		username = :username,
		refresh_started_at = :refresh_started_at,
		refresh_active_at = :refresh_active_at,
		refresh_completed_at = :refresh_completed_at,
		refresh_status = :refresh_status
	WHERE id = :id`

	_, err := s.db.NamedExecContext(ctx, q, account)
	return err
}
