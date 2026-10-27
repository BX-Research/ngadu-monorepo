package dal

import (
	"context"

	"github.com/BX-Research/ngadu-monorepo/backend/app/types"
	"github.com/jackc/pgx/v5"
)

// Querier is the database capability required by user lookups.
type Querier interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

type Users struct {
	db Querier
}

func NewUsers(db Querier) *Users {
	return &Users{db: db}
}

func (u *Users) FindUser(ctx context.Context, id string) (types.UserInfo, error) {
	var user types.UserInfo
	err := u.db.QueryRow(ctx, `SELECT id::text, display_name FROM users WHERE id = $1`, id).
		Scan(&user.ID, &user.DisplayName)
	return user, err
}
