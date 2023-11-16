package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/khhini/riset-gcp-api-gateway-auth/internal/modules/user/entity"
)

const (
	insertUserQuery        = `INSERT INTO users(id, username, email, password, created_on) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	getUserByIDQuery       = `SELECT id, username, email, created_on, last_login FROM users WHERE id = $1`
	getUserByUsernameQuery = `SELECT id, username, password, email, created_on, last_login FROM users WHERE username = $1`
	updateUserQuery        = `UPDATE users SET username = $2, email = $3, last_login = $4 WHERE id = $1`
	deleteUserQuery        = `DELETE FROM users WHERE id = $1`
)

type PostgreRepository struct {
	ConnPool *pgxpool.Pool
}

func NewPostgreRepository(connPool *pgxpool.Pool) entity.Repository {
	return &PostgreRepository{
		ConnPool: connPool,
	}
}

func (r *PostgreRepository) InsertUser(ctx context.Context, user *entity.User) (userID uuid.UUID, err error) {

	err = r.ConnPool.QueryRow(
		ctx,
		insertUserQuery,
		uuid.New(),
		user.Username,
		user.Email,
		user.Password,
		user.CreatedOn.Time,
	).Scan(&userID)

	if err != nil {
		return uuid.Nil, fmt.Errorf("error inserting user: %w", err)
	}
	return userID, nil
}

func (r *PostgreRepository) GetUserByID(ctx context.Context, userID string) (*entity.User, error) {
	var user entity.User

	row := r.ConnPool.QueryRow(ctx, getUserByIDQuery, userID)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedOn, &user.LastLogin)

	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("error fetching user: %w", err)
	}

	return &user, nil
}

func (r *PostgreRepository) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User

	row := r.ConnPool.QueryRow(ctx, getUserByUsernameQuery, username)
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.CreatedOn, &user.LastLogin)

	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("error fetching user: %w", err)
	}

	return &user, nil
}

func (r *PostgreRepository) UpdateUser(ctx context.Context, userID string, updatedUser *entity.User) (err error) {
	_, err = r.ConnPool.Exec(
		ctx,
		updateUserQuery,
		userID,
		updatedUser.Username,
		updatedUser.Email,
		updatedUser.LastLogin,
	)

	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}
	return nil
}

func (r *PostgreRepository) DeleteUser(ctx context.Context, userID string) (err error) {
	_, err = r.ConnPool.Exec(ctx, deleteUserQuery, userID)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	return nil
}
