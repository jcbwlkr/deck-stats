package auth

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/jcbwlkr/deck-stats/internal/domains/users"
)

type Authenticator struct {
	Secret string
}

func NewAuthenticator(secret string) *Authenticator {
	return &Authenticator{
		Secret: secret,
	}
}

type Claims struct {
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

func (a *Authenticator) GenerateJWT(user *users.User) (string, error) {
	var claims Claims
	claims.Subject = user.ID
	claims.Username = user.Username
	claims.Roles = user.Roles
	claims.IssuedAt = jwt.NewNumericDate(time.Now())
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.Secret))
}

func (a *Authenticator) ValidateToken(t string) (users.User, error) {

	var claims Claims

	token, err := jwt.ParseWithClaims(t, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(a.Secret), nil
	})
	if err != nil {
		return users.User{}, err
	}

	if !token.Valid {
		return users.User{}, err
	}

	return users.User{
		ID:       claims.Subject,
		Username: claims.Username,
		Roles:    claims.Roles,
	}, nil
}

func (a *Authenticator) Middleware() func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			parts := strings.Split(r.Header.Get("Authorization"), " ")
			if len(parts) != 2 {
				http.Error(w, "missing or malformed authorization header", http.StatusUnauthorized)
				return
			}

			user, err := a.ValidateToken(parts[1])
			if err != nil {
				slog.WarnContext(ctx, "invalid authentication", "error", err)
				http.Error(w, "invalid authentication", http.StatusUnauthorized)
				return
			}

			ctx = StoreUser(ctx, user)
			r = r.WithContext(ctx)

			next(w, r)
		})
	}
}
