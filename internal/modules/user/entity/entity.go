package entity

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
)

type User struct {
	SafeUser
	Password string `json:"password"`
}

type SafeUser struct {
	ID        uuid.UUID    `json:"id"`
	Username  string       `json:"username"`
	Email     string       `json:"email"`
	CreatedOn sql.NullTime `json:"created_on"`
	LastLogin sql.NullTime `json:"last_login"`
}

func (u *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(SafeUser{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedOn: u.CreatedOn,
		LastLogin: u.LastLogin,
	})
}

type Repository interface {
	InsertUser(ctx context.Context, user *User) (userID uuid.UUID, err error)
	GetUserByID(ctx context.Context, userID string) (user *User, err error)
	GetUserByUsername(ctx context.Context, username string) (user *User, err error)
	UpdateUser(ctx context.Context, userID string, updatedUser *User) (err error)
	DeleteUser(ctx context.Context, userID string) (err error)
}

type Service interface {
	InsertUser(ctx context.Context, data *User) (userID uuid.UUID, err error)
	GetUserByID(ctx context.Context, userID string) (user *User, err error)
	UpdateUser(ctx context.Context, userID string, updatedUser *User) (err error)
	DeleteUser(ctx context.Context, userID string) (err error)
}
