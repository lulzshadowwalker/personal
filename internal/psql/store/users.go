package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lulzshadowwalker/personal/internal"
	"github.com/lulzshadowwalker/personal/internal/psql/db"
)

var ErrUserNotFound = errors.New("user not found")

type UserStore struct {
	q *db.Queries
}

func NewUserStore(db *db.Queries) *UserStore {
	return &UserStore{
		q: db,
	}
}

func (s *UserStore) toEntity(user db.User) internal.User {
	return internal.User{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
	}
}

func (s *UserStore) GetUserByEmail(ctx context.Context, email string) (internal.User, error) {
	user, err := s.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return internal.User{}, ErrUserNotFound
		}

		return internal.User{}, err
	}

	return s.toEntity(user), nil
}

func (s *UserStore) CreateUser(ctx context.Context, user internal.User) (internal.User, error) {
	userDB, err := s.q.CreateUser(ctx, db.CreateUserParams{
		Name:  user.Name,
		Email: user.Email,
	})
	if err != nil {
		return internal.User{}, err
	}

	return s.toEntity(userDB), nil
}
