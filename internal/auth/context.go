package auth

import (
	"context"

	"github.com/jcbwlkr/deck-stats/internal/domains/users"
)

type userKeyT string

const userKey = userKeyT("user")

func StoreUser(ctx context.Context, user users.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) (users.User, bool) {
	u, ok := ctx.Value(userKey).(users.User)
	return u, ok
}
