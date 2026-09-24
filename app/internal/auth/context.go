package auth

import (
	"context"
	"scavenger/internal/domain"
)

type ctxKey int

const userKey ctxKey = 1

func WithUser(ctx context.Context, u *domain.User) context.Context {
	return context.WithValue(ctx, userKey, u)
}

func UserFromCtx(ctx context.Context) *domain.User {
	u, _ := ctx.Value(userKey).(*domain.User)
	return u
}
